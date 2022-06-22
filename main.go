package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mizumoto-cn/gobalancer/proxy"
)

func main() {
	//load and check config
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	//PrintBanner()
	err = config.Validate()
	if err != nil {
		log.Fatal("Failed to validate config:", err)
	}

	// create and start proxy
	router := mux.NewRouter()
	for _, location := range config.Location {
		httpProxy, err := proxy.NewHttpProxy(location.ProxyPass, location.BalanceMode)
		if err != nil {
			log.Fatal("Failed to create proxy:", err)
		}
		// health check
		if config.HealthCheck {
			httpProxy.HealthCheck(config.HealthCheckInterval)
		}
		// register proxy with each location
		router.Handle(location.Pattern, httpProxy)
	}

	// setup with config
	if config.MaxConnections > 0 {
		router.Use(maxConnectionsMiddleware(config.MaxConnections))
	}
	svr := &http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: router,
	}
	// print config
	PrintConfig(config)

	// listen and serve
	switch config.Schema {
	case "http":
		err = svr.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to listen and serve:", err)
		}
	case "https":
		err = svr.ListenAndServeTLS(config.SSLCert, config.SSLCertKey)
		if err != nil {
			log.Fatal("Failed to listen and serve:", err)
		}
	}
}

// maxConnectionsMiddleware is a middleware that limits the number of connections
// to a given server.
func maxConnectionsMiddleware(maxConnections uint64) mux.MiddlewareFunc {
	sem := make(chan struct{}, maxConnections) // semaphore
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acquire()
			defer release()
			next.ServeHTTP(w, r)
		})
	}
}

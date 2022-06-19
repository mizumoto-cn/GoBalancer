package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/mizumoto-cn/gobalancer/balancer"
)

// I hate arbitrary names in http protocol.
// For god's sake, what does the fucking letter "X" mean?
// And everyone using these names, fucking moron.
// See why they choose dumb-ass "x-" prefixes at
// https://datatracker.ietf.org/doc/html/rfc6648
var (
	XRealIP       = http.CanonicalHeaderKey("X-Real-IP")
	XProxyIP      = http.CanonicalHeaderKey("X-Proxy-IP")
	XForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
)

var (
	ReverseProxy = "Balancer-ReverseProxy"
)

type HTTPProxy struct {
	hostMap  map[string]*httputil.ReverseProxy
	balancer balancer.Balancer

	sync.RWMutex // protects maps as they are shared across goroutines
	alive        map[string]bool
}

// NewHTTPProxy creates a new HTTPProxy instance.
// targets is a slice of URLs that the proxy will forward requests to.
// balancerAlgorithm is the name of the balancer algorithm.
// This function resolves every URL it receives and creates a new reverse proxy for each.
func NewHttpProxy(targets []string, balancerAlgorithm string) (*HTTPProxy, error) {
	hostMap := make(map[string]*httputil.ReverseProxy)
	hosts := make([]string, 0)
	alive := make(map[string]bool)

	for _, target := range targets {
		url, err := url.Parse(target)
		if err != nil {
			return nil, err
		}
		proxy := httputil.NewSingleHostReverseProxy(url)

		originDirect := proxy.Director
		// proxy.Director is a function that takes a request and returns a modified request.
		proxy.Director = func(r *http.Request) {
			originDirect(r)
			r.Header.Set(XProxyIP, ReverseProxy) // set the proxy IP to distinguish from real IP
			r.Header.Set(XRealIP, GetIP(r))      // set the real IP
		}

		host := GetHost(url) // get the host name from the URL
		alive[host] = true   // set the host as alive as default
		hostMap[host] = proxy
		hosts = append(hosts, host)
	}

	balancer, err := balancer.Build(balancerAlgorithm, hosts)
	if err != nil {
		return nil, err
	}

	return &HTTPProxy{
		hostMap:  hostMap,
		balancer: balancer,
		alive:    alive,
	}, nil
}

// To implement reverse proxy, we need to override ServeHTTP method.
// HTTPProxy will forward the request to the target host based on the balancer algorithm.
// When failed to forward the request, HTTPProxy will return a 502 error.
func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("proxy panic: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(err.(error).Error()))
		}
	}()

	host, err := h.balancer.BalanceHost(GetIP(r))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("balancer error: %v", err)))
		return
	}

	h.balancer.Inc(host)
	defer h.balancer.Done(host)
	h.hostMap[host].ServeHTTP(w, r)
}

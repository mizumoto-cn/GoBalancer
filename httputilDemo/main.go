package main

//
// @title go-httputil-demo
// A simple example of using the go-httputil package as a reverse proxy.

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HTTPProxy struct {
	proxy *httputil.ReverseProxy
}

func main() {
	proxy, err := NewHTTPProxy("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	http.Handle("/", proxy)
	http.ListenAndServe(":8080", nil)
}

func NewHTTPProxy(target string) (*HTTPProxy, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	return &HTTPProxy{
		proxy: httputil.NewSingleHostReverseProxy(url),
	}, nil
}

func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(w, r)
}

// In the code above, we defined a simple HTTP proxy that forwards requests to a
// remote server. HTTPProxy is a struct containing a ReverseProxy instance.
// We resolved the target URL into a *url.URL instance and passed it to the
// httputil.NewSingleHostReverseProxy function. This function returns a
// ReverseProxy instance that will forward requests to the remote server.

// HTTPProxy needs to implement the ServeHTTP method.

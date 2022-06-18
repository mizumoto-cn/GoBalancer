package proxy

import (
	"net/http/httputil"
	"sync"

	"github.com/mizumoto-cn/gobalancer/balancer"
)

type HTTPProxy struct {
	hostMap  map[string]*httputil.ReverseProxy
	balancer balancer.Balancer

	sync.RWMutex // protects maps as they are shared across goroutines
	alive        map[string]bool
}

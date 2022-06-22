package proxy

import (
	"log"
	"time"
)

// GetServiceAliveStatus reads the alive status of the service.
// service: url format
func (p *HTTPProxy) GetServiceAliveStatus(service string) bool {
	// Lock as map is not thread safe.
	p.RLock()
	defer p.RUnlock()
	return p.alive[service]
}

// SetServiceAliveStatus sets the alive status of the service.
func (p *HTTPProxy) SetServiceAliveStatus(service string, alive bool) {
	// Lock as map is not thread safe.
	p.Lock()
	defer p.Unlock()
	p.alive[service] = alive
}

// HealthCheck goroutine checks the health of the service.
func (p *HTTPProxy) HealthCheck(interval uint64) {
	// for each host in the pool execute a health check goroutine.
	for host := range p.hostMap {
		go p.healthCheck(host, interval)
	}
}

// healthCheck
func (p *HTTPProxy) healthCheck(host string, interval uint64) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for range ticker.C {
		if p.GetServiceAliveStatus(host) && !IsBackendAlive(host) {
			log.Printf("Service %s is down, removing from the pool.", host)
			p.SetServiceAliveStatus(host, false)
			p.balancer.RemoveHost(host)
		} else if !p.GetServiceAliveStatus(host) && IsBackendAlive(host) {
			log.Printf("Service %s is up, adding to the pool.", host)
			p.SetServiceAliveStatus(host, true)
			p.balancer.AddHost(host)
		}
	}
}

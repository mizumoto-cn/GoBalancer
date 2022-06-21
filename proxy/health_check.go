package proxy

import (
	"log"
	"time"
)

// GetServiceAliveStatus reads the alive status of the service.
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
func (p *HTTPProxy) HealthCheck( interval int64){
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for range ticker.C {
		
}
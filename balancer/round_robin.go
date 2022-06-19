package balancer

import (
	"sync"

	"github.com/mizumoto-cn/gobalancer/util"
)

// RoundRobin selects hosts in a round robin fashion.
// In other words, it selects the next host in the list.
type RoundRobin struct {
	sync.RWMutex
	i     uint64 // index of the next host to be selected
	hosts []string
}

// Register RoundRobin as a balancer algorithm in the factories map.
func init() {
	factoriesMap["round_robin"] = NewRoundRobin // register the factory function
}

func NewRoundRobin(hosts []string) (Balancer, error) {
	return &RoundRobin{i: 0, hosts: hosts}, nil
}

// AddHost adds a host to the list of hosts.
func (b *RoundRobin) AddHost(host string) error {
	b.Lock()
	defer b.Unlock()
	if !util.Contains(b.hosts, host) {
		b.hosts = append(b.hosts, host)
	}
	return nil
}

// RemoveHost removes a host from the list of hosts.
func (b *RoundRobin) RemoveHost(host string) error {
	b.Lock()
	defer b.Unlock()
	if util.Contains(b.hosts, host) {
		b.hosts = util.Remove(b.hosts, host)
	}
	return nil
}

// BalanceHost returns the next host in the list.
func (b *RoundRobin) BalanceHost(key string) (string, error) {
	b.RLock()
	defer b.RUnlock()
	if len(b.hosts) == 0 {
		return "", ErrHostNotFound
	}
	host := b.hosts[b.i%uint64(len(b.hosts))]
	b.i++
	return host, nil
}

// Inc increases the number of connections to the host by 1.
func (b *RoundRobin) Inc(host string) error {
	b.Lock()
	defer b.Unlock()
	return nil
}

// Done decreases the number of connections to the host by 1.
func (b *RoundRobin) Done(host string) error {
	b.Lock()
	defer b.Unlock()
	return nil
}

package balancer

// ip_hash implements a IP hash balancer.

import (
	"hash/crc32"
	"sync"

	"github.com/mizumoto-cn/gobalancer/util"
)

// registration of balancer algorithms
func init() {
	factoriesMap["ip_hash"] = NewIPHash
}

// IPHash will select a host based on the IP address of the client.
type IPHash struct {
	sync.RWMutex
	hosts []string // list of hosts
}

// NewIPHash returns a IPHash Balancer instance.
func NewIPHash(hosts []string) (Balancer, error) {
	return &IPHash{hosts: hosts}, nil
}

// AddHost adds a host to the list of hosts.
func (b *IPHash) AddHost(host string) error {
	b.Lock()
	defer b.Unlock()
	if !util.Contains(b.hosts, host) {
		b.hosts = append(b.hosts, host)
	}
	return nil
}

// RemoveHost removes a host from the list of hosts.
func (b *IPHash) RemoveHost(host string) error {
	b.Lock()
	defer b.Unlock()
	if util.Contains(b.hosts, host) {
		b.hosts = util.Remove(b.hosts, host)
	}
	return nil
}

// BalanceHost returns a host based on the IP address of the client.
func (b *IPHash) BalanceHost(key string) (string, error) {
	b.RLock()
	defer b.RUnlock()
	if len(b.hosts) == 0 {
		return "", ErrHostNotFound
	}
	// calculate the hash of the client IP address and modulus the number of hosts
	hashedIP := crc32.ChecksumIEEE([]byte(key)) % uint32(len(b.hosts))
	return b.hosts[hashedIP], nil
}

// Inc increases the number of connections to the host by 1.
func (b *IPHash) Inc(host string) error {
	return nil
}

// Dec decreases the number of connections to the host by 1.
func (b *IPHash) Done(host string) error {
	return nil
}

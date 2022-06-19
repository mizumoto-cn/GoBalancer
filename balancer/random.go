package balancer

import (
	"math/rand"
	"sync"
	"time"

	"github.com/mizumoto-cn/gobalancer/util"
)

// registration of balancer algorithms
func init() {
	factoriesMap["random"] = NewRandom
}

// Random selects hosts in a random fashion.
type Random struct {
	sync.RWMutex
	hosts []string // list of hosts
	rand  *rand.Rand
}

// NewRandom returns a Random Balancer instance.
func NewRandom(hosts []string) (Balancer, error) {
	return &Random{hosts: hosts, rand: rand.New(rand.NewSource(time.Now().UnixNano()))}, nil
}

// AddHost adds a host to the list of hosts.
func (b *Random) AddHost(host string) error {
	b.Lock()
	defer b.Unlock()
	if !util.Contains(b.hosts, host) {
		b.hosts = append(b.hosts, host)
	}
	return nil
}

// RemoveHost removes a host from the list of hosts.
func (b *Random) RemoveHost(host string) error {
	b.Lock()
	defer b.Unlock()
	if util.Contains(b.hosts, host) {
		b.hosts = util.Remove(b.hosts, host)
	}
	return nil
}

// BalanceHost returns a random host from the list.
func (b *Random) BalanceHost(key string) (string, error) {
	b.RLock()
	defer b.RUnlock()
	if len(b.hosts) == 0 {
		return "", ErrHostNotFound
	}
	host := b.hosts[b.rand.Intn(len(b.hosts))]
	return host, nil
}

// Inc increases the number of connections to the host by 1.
func (b *Random) Inc(host string) error {
	return nil
}

// Done decreases the number of connections to the host by 1.
func (b *Random) Done(host string) error {
	return nil
}

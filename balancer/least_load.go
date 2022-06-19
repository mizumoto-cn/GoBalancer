package balancer

import (
	"sync"

	fibHeap "github.com/mizumoto-cn/gobalancer/util"
)

// LeastLoad is a balancer that selects a host based on the least load.
// It is a fibonacci heap. See https://en.wikipedia.org/wiki/Fibonacci_heap
// for more information. It is a min-heap, which takes O(1) time to find the
// minimum element, which is the host with the least load here. It takes O(log(n))
// time to insert and delete an element.

// registration of least load balancer algorithm
func init() {
	factoriesMap["least_load"] = NewLeastLoad
}

// Tag is a tag for LeastLoad.
func (h *host_info) Tag() interface{} {
	return h.host_name
}

// Key is a key for LeastLoad.
func (h *host_info) Key() float64 {
	return float64(h.load)
}

// LeastLoad will choose a host based on the least load.
type LeastLoad struct {
	sync.RWMutex
	heap *fibHeap.FibHeap
}

// NewLeastLoad creates a new LeastLoad instance.
func NewLeastLoad(hosts []string) (Balancer, error) {
	ll := &LeastLoad{
		heap: fibHeap.NewFibHeap(),
	}
	for _, host := range hosts {
		ll.AddHost(host)
	}
	return ll, nil
}

// AddHost adds a host to the LeastLoad instance.
func (ll *LeastLoad) AddHost(host string) error {
	ll.Lock()
	defer ll.Unlock()
	if ok := ll.heap.GetValue(host); ok == nil {
		return ErrHostAlreadyExists
	}
	// h := &host_info{
	// 	host_name: host,
	// 	load:      0,
	// }
	ll.heap.Insert(host, 0)
	return nil
}

// RemoveHost removes a host from the LeastLoad instance.
func (ll *LeastLoad) RemoveHost(host string) error {
	ll.Lock()
	defer ll.Unlock()
	if ok := ll.heap.GetValue(host); ok == nil {
		return ErrHostNotFound
	}
	ll.heap.Delete(host)
	return nil
}

// BalanceHost chooses a host based on the least load.
func (ll *LeastLoad) BalanceHost(tag string) (string, error) {
	ll.RLock()
	defer ll.RUnlock()
	if ll.heap.Num() == 0 {
		return "", ErrHostNotFound
	}
	return ll.heap.MinimumValue().Tag().(string), nil
}

// Inc increases the number of connections to a host by 1.
func (ll *LeastLoad) Inc(host string) error {
	ll.Lock()
	defer ll.Unlock()
	if ok := ll.heap.GetValue(host); ok == nil {
		return ErrHostNotFound
	}
	h := ll.heap.GetValue(host).(*host_info)
	h.load++
	ll.heap.IncreaseKeyValue(h)
	return nil
}

// Done decreases the number of connections to a host by 1.
func (ll *LeastLoad) Done(host string) error {
	ll.Lock()
	defer ll.Unlock()
	if ok := ll.heap.GetValue(host); ok == nil {
		return ErrHostNotFound
	}
	h := ll.heap.GetValue(host).(*host_info)
	h.load--
	ll.heap.DecreaseKeyValue(h)

	return nil
}

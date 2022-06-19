package balancer

// p2c balancer algorithm
import (
	"math/rand"
	"sync"
	"time"

	"github.com/mizumoto-cn/gobalancer/util"
)

const Salt = "mizu#salt"

// registration of power of two choices balancer algorithm
func init() {
	factoriesMap["power_of_two_choices"] = NewPowerOfTwoChoices
}

type host_info struct {
	host_name string // host name
	load      uint64 // load of the host
}

// PowerOfTwoChoices will select a host based on the IP address of the client.
type PowerOfTwoChoices struct {
	sync.RWMutex
	hosts   []*host_info // hosts
	rand    *rand.Rand
	loadMap map[string]*host_info // snapshot of hosts
}

// NewPowerOfTwoChoices creates a new PowerOfTwoChoices instance.
func NewPowerOfTwoChoices(hosts []string) (Balancer, error) {
	if len(hosts) == 0 {
		return nil, ErrHostNotFound
	}
	p := &PowerOfTwoChoices{
		hosts:   []*host_info{},
		loadMap: make(map[string]*host_info),
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	for _, host := range hosts {
		p.AddHost(host)
	}
	return p, nil
}

// AddHost adds a host to the PowerOfTwoChoices instance.
func (p *PowerOfTwoChoices) AddHost(host string) error {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.loadMap[host]; ok {
		return ErrHostAlreadyExists
	}
	h := &host_info{
		host_name: host,
		load:      0,
	}
	p.hosts = append(p.hosts, h)
	p.loadMap[host] = h

	return nil
}

//delete_host_info deletes a host_info from slice
func (p *PowerOfTwoChoices) delete_host_info(host string) {
	for i, v := range p.hosts {
		if v.host_name == host {
			p.hosts = append(p.hosts[:i], p.hosts[i+1:]...)
			return
		}
	}
}

// RemoveHost removes a host from the PowerOfTwoChoices instance.
func (p *PowerOfTwoChoices) RemoveHost(host string) error {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.loadMap[host]; !ok {
		return ErrHostNotFound
	}
	delete(p.loadMap, host)
	p.delete_host_info(host)
	return nil
}

// BalanceHost selects a host based on the IP address of the client.
func (p *PowerOfTwoChoices) BalanceHost(clientIP string) (string, error) {
	p.RLock()
	defer p.RUnlock()
	if len(p.hosts) == 0 {
		return "", ErrHostNotFound
	}
	k1, k2 := p.hash(clientIP)
	n1 := p.hosts[k1%uint32(len(p.hosts))].host_name
	n2 := p.hosts[k2%uint32(len(p.hosts))].host_name

	host := n2

	// choose a less loaded host
	if p.loadMap[n1].load < p.loadMap[n2].load {
		host = n1
	}
	return host, nil
}

func (p *PowerOfTwoChoices) hash(clientIP string) (uint32, uint32) {
	var k1, k2 uint32
	if clientIP == "" {
		k1 = uint32(p.rand.Intn(len(p.hosts)))
		k2 = uint32(p.rand.Intn(len(p.hosts)))
	} else {
		k1, k2 = util.Hash(clientIP, Salt, len(p.hosts))
	}
	return k1, k2
}

// Inc increases the number of connections to the host by 1.
func (p *PowerOfTwoChoices) Inc(host string) error {
	p.Lock()
	defer p.Unlock()

	h, ok := p.loadMap[host]

	if !ok {
		return ErrHostNotFound
	}
	h.load++
	return nil
}

// Done decreases the number of connections to the host by 1.
func (p *PowerOfTwoChoices) Done(host string) error {
	p.Lock()
	defer p.Unlock()

	h, ok := p.loadMap[host]

	if !ok {

		return ErrHostNotFound

	}
	if h.load > 0 {
		h.load--
		return nil
	} else {
		return ErrNoLoadToRemove
	}
}

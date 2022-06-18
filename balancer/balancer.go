package balancer

import "errors"

type Balancer interface {
	AddHost(host string) error
	RemoveHost(host string) error
	BalanceHost(key string) (string, error) // choose a host from the list regarding the key
	Inc(host string) error                  // increase the number of connections to the host by 1
	Done(host string) error                 // decrease the number of connections to the host by 1
}

var (
	ErrHostNotFound              = errors.New("host not found")
	ErrBalancerAlgorithmNotFound = errors.New("balancer algorithm not found")
)

type factories func([]string) (Balancer, error)

var factoriesMap = map[string]factories

func Build (algorithm string, hosts []string) (Balancer, error) {
	if factory, ok := factoriesMap[algorithm]; ok {
		return factory(hosts), nil
	}
	return nil, ErrBalancerAlgorithmNotFound
}
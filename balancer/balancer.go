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
	ErrHostAlreadyExists         = errors.New("host already exists")
	ErrNoLoadToRemove            = errors.New("no load to remove")
)

// factory design pattern
type factory func([]string) (Balancer, error)

// factoriesMap is a map of algorithm name to factory function
var factoriesMap = make(map[string]factory)

// Build returns a Balancer instance based on the algorithm name
func Build(algorithm string, hosts []string) (Balancer, error) {
	if factory, ok := factoriesMap[algorithm]; ok {
		return factory(hosts) //, nil
	}
	return nil, ErrBalancerAlgorithmNotFound
}

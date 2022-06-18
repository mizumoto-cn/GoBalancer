package balancer

type Balancer interface {
	AddHost(host string) error
	RemoveHost(host string) error
	BalanceHost(key string) (string, error)
	Inc(host string) error
	Done(host string) error
}

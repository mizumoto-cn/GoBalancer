# Go Balancer

![Go Version](https://img.shields.io/badge/Go%20version-v1.18.3-yellow)

A tiny payload balancer written in go. Based on net/http/httputil. An layer-7 application and also a payload-balancing algorithm library.

## Quick Start

nil

## Forward Proxy and Reverse Proxy

Maybe most of us have used proxies. A forward proxy can be roughly defined as a domain between user client and the internet, which hides the real client from the internet, acting as if they(proxies) are users themselves.

And vice versa. So there are reverse proxies hiding concrete servers from the client. Like the pic shown beneath.

![Forward and Reverse](arc/reverse_proxy.png)

So here in [httputilDemo](httputilDemo/main.go), I will show you how to create a simple reverse proxy.

## Architecture

Try show the architecture of the project in tree diagram below.

```boo
.
├── balancer # Load balancers
├── proxy    # proxies
└── util     # utilities
```

Go Balancer is a light-weight payload balancer.
It has no complex architecture, basically only uses the factory pattern in balancer registry and creation.

Let's start with the payload balancer part.

### Balancer

Balancer is a interface that defines the payload balancer with the following methods.

[Balancer](balancer/balancer.go) <!-- markdownlint-disable MD010-->

```golang
type Balancer interface {
	AddHost(host string) error
	RemoveHost(host string) error

     // choose a host from the list regarding the key
	BalanceHost(key string) (string, error)

    // increase the number of connections to the host by 1
	Inc(host string) error          

    // decrease the number of connections to the host by 1        
	Done(host string) error                 
}
```

After understanding the abstract of Balancer, let's start to implement the balancer algorithms.

There will be 7 algorithms implemented in this project:

- random
- round-robin
- power-of-two random choice
- consistent hash
- consistent hash with bounded capacity
- ip-hash
- least-loaded

### Factory Pattern

The factory pattern is used to create the balancer. We defined a `Factory` function that returns a `Balancer` interface.

Then we use `Build()` function to create the balancer through the factory by calling the `factory` function.

The factory pattern is defined as follows at [balancer.go](balancer/balancer.go).

```golang
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
```

For each load balancer algorithm, we register them into the `factories` map in `init()` function.

Take [`round-robin`](balancer/round_robin.go) as an example:

```golang
// Register RoundRobin as a balancer algorithm in the factories map.
func init() {
	factoriesMap["round_robin"] = NewRoundRobin // register the factory function
}

func NewRoundRobin(hosts []string) (Balancer, error) {
	return &RoundRobin{i: 0, hosts: hosts}, nil
}
```

### Health Check

> To be implemented.

### Util & Configure

We set separated utilities in the [`util` package](util/util.go), while the configuration is in the `main` package.

> To be implemented.

## Algorithms

Typically, we have a set of 7 load balancer algorithms:

- random [![status](https://img.shields.io/badge/status-implemented-green)](/balancer/random.go)
- round-robin [![status](https://img.shields.io/badge/status-implemented-green)](/balancer/round_robin.go)
- power-of-two random choice [![status](https://img.shields.io/badge/status-implemented-green)](/balancer/p2c.go)
- consistent ![status](https://img.shields.io/badge/status-not%20implemented-gray)
- consistent hash with bounded capacity ![status](https://img.shields.io/badge/status-not%20implemented-gray)
- ip-hash [![status](https://img.shields.io/badge/status-implemented-green)](/balancer/ip_hash.go)
- least-loaded [![status](https://img.shields.io/badge/status-implemented-green)](/balancer/least_loaded.go)
  
Thank tencentyun/tsf-go for [practical p2c algorithm implementation](https://github.com/tencentyun/tsf-go/blob/master/balancer/p2c/p2c.go) examples.

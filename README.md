# Go Balancer

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

Go Balancer is a payload balancer, so let's start with the payload balancer part.

### Balancer

Balancer is a interface that defines the payload balancer with the following methods.

[Banlancer](balancer/balancer.go)

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
# Go Balancer

![Go Version](https://img.shields.io/badge/Go%20version-v1.18.3-yellow)
[![License](https://img.shields.io/badge/License-MGPL%20v1.2-green.svg)](/License/Mizumoto%20General%20Public%20License%20v1.2.md)

[![Build](https://github.com/mizumoto-cn/GoBalancer/actions/workflows/go.yml/badge.svg)](https://github.com/mizumoto-cn/GoBalancer/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/mizumoto-cn/gobalancer)](https://goreportcard.com/report/github.com/mizumoto-cn/gobalancer)
[![CodeFactor](https://www.codefactor.io/repository/github/mizumoto-cn/gobalancer/badge)](https://www.codefactor.io/repository/github/mizumoto-cn/gobalancer)
[![codecov](https://codecov.io/gh/mizumoto-cn/GoBalancer/branch/main/graph/badge.svg?token=UK8P1TX6WE)](https://codecov.io/gh/mizumoto-cn/GoBalancer)
[![Go Version](https://img.shields.io/badge/Go%20version-1.18.3-green.svg)](https://golang.org/doc/install)

```golang
  ________      __________        .__                                    
 /  _____/  ____\______   \_____  |  | _____    ____   ____  ___________ 
/   \  ___ /  _ \|    |  _/\__  \ |  | \__  \  /    \_/ ___\/ __ \_  __ \
\    \_\  (  <_> )    |   \ / __ \|  |__/ __ \|   |  \  \__\  ___/|  | \/
 \______  /\____/|______  /(____  /____(____  /___|  /\___  >___  >__|   
        \/              \/      \/          \/     \/     \/    \/
```  

🚦 A fast, stable, light-weight layer-7 load balancer written in go. Based on `net/http/httputil`, also a load-balancing algorithm library.

## Quick Start

### Install

First, clone the repository,

```bash
  git clone https://github.com/mizumoto-cn/GoBalancer.git
```

and then build the binary with:

```bash
  cd GoBalancer
  go build
```

Or simply fetch the docker image:

```bash
  docker pull mizumotocn/go-balancer
```

### First Run

You need to create a config file first, see example config file at [config.json](./config.json).

And then, you can run the binary.

On Linux / UNIX, you can type `./GoBalancer` to run the binary.

```bash
  ./GoBalancer
```

On Windows, you can use `GoBalancer.exe` to run the binary.

```bash
  ./GoBalancer.exe
```

And you'll see the following output:

```bash
PS D:\Reposits\GoBalancer> .\gobalancer.exe
Schema: http
Port: 8089
HealthCheck: true
MaxConnections: 100

  ________      __________        .__
 /  _____/  ____\______   \_____  |  | _____    ____   ____  ___________
/   \  ___ /  _ \|    |  _/\__  \ |  | \__  \  /    \_/ ___\/ __ \_  __ \
\    \_\  (  <_> )    |   \ / __ \|  |__/ __ \|   |  \  \__\  ___/|  | \/
 \______  /\____/|______  /(____  /____(____  /___|  /\___  >___  >__|
        \/              \/      \/          \/     \/     \/    \/

Pattern: /
ProxyPass: [http://192.168.1.1 http://192.168.1.2:1015 https://192.168.1.2 http://my-server.com]
BalanceMode: round_robin
```

The `GoBalancer` will perform periodical `health check`s on all proxy sites. When a site is unreachable, it will be automatically removed from the balancer, which will look like:

```bash
2022/06/22 18:48:21 Service 192.168.1.2:1015 is down, removing from the pool.
2022/06/22 18:48:21 Service 192.168.1.1:80 is down, removing from the pool.
2022/06/22 18:48:21 Service 192.168.1.2:443 is down, removing from the pool.
```

However, the balancer will still perform health checks on unreachable sites. When a site is reachable, it will automatically be added to the balancer.

You can also disable the health check by setting `"health_check": false` in the config file.

For more complicated architecture scenarios, I'd prefer to change the health check mechanism, making the sites send heartbeat messages to the balancer to get registered. And sites shall be wiped out from the pool when the heartbeat message is not received for a certain number of heartbeat periods.

### API Reference

GoBalancer is also a go library that implements load-balancing algorithms, it can be used as API on its own.

What you need to do is to import the `gobalancer` package first.

```bash
go get github.com/mizumoto-cn/gobalancer
```

Then

```go
import "github.com/mizumoto-cn/gobalancer"
```

Now you can use `balancer.Build` to build a balancer.

<!-- markdownlint-disable no-hard-tabs -->

```go
hosts := []string{
	"http://192.168.97.1",
	"http://192.168.97.2",
	"http://192.168.97.3",
	"http://192.168.97.4",
}

balancer, err := balancer.Build(balancer.PowerOfTwoChoicesBalancer, hosts)
if err != nil {
	return err
}
```

Balancers were built with the factory pattern, and you can find the aliases in [/balancer/const.go](/balancer/const.go). You can find the implemented balancers in [Algorithms](#algorithms) part. Balancers implement the `Balancer` interface:

```go
type Balancer interface {
	AddHost(host string) error
	RemoveHost(host string) error
	BalanceHost(key string) (string, error) // choose a host from the list regarding the key
	Inc(host string) error                  // increase the number of connections to the host by 1
	Done(host string) error                 // decrease the number of connections to the host by 1
}
```

Balancers shall be used like this:

```go
client := "http://192.168.97.11" // request client ip
target, err := balancer.Balance(client)
if err != nil {
	return err
}

balancer.Inc(target)// increase the connection count of the target server
defer balancer.Done(target)// decrease the connection count of the target server

// route to the target server, do something
...
```

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

### Forward Proxy and Reverse Proxy

Maybe most of us have used proxies. A forward proxy can be roughly defined as a domain between user client and the internet, which hides the real client from the internet, acting as if they(proxies) are users themselves.

And vice versa. So there are reverse proxies hiding concrete servers from the client. Like the pic shown beneath.

![Forward and Reverse](arc/reverse_proxy.png)

So here in [httputilDemo](httputilDemo/main.go), I will show you how to create a simple reverse proxy.

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

In general, we follow these rules to set health check tests in multi-service architecture:

- between service providers and consumers, making sure the service invocation is available
- between registrar and service providers, making sure the registry knows the service provider is available
- between registrar and service consumers, making sure the registry knows the service consumer is available

Herein, both service providers and consumers are the clients to the registrar.

We don't really need to implement this project to be part of a multi-service architecture thus we only set health check tests between the proxy and the service-providing servers( as clients).

I'd rather prefer clients sending heartbeats to the registrar regularly than to implement the health check tests in the proxy but it will make the project too complicated for its potential application scenarios.

### Util & Configure

We set separated utilities in the [`util` package](util/util.go), while the configuration is in the `main` package.

> To be implemented.

### Algorithms

We have a set of balancer algorithms:

- random [![status](https://img.shields.io/badge/test-pass-green)](/balancer/random.go)
- round-robin [![status](https://img.shields.io/badge/test-pass-green)](/balancer/round_robin.go)
- power-of-two random choice [![status](https://img.shields.io/badge/test-pass-green)](/balancer/p2c.go)
- consistent ![status](https://img.shields.io/badge/status-not%20implemented-red)
- consistent hash with bounded capacity ![status](https://img.shields.io/badge/status-not%20implemented-red)
- ip-hash [![status](https://img.shields.io/badge/test-pass-green)](/balancer/ip_hash.go)
- least-loaded [![status](https://img.shields.io/badge/test-pass-green)](/balancer/least_loaded.go)
  
Thank tencentyun/tsf-go for [practical p2c algorithm implementation](https://github.com/tencentyun/tsf-go/blob/master/balancer/p2c/p2c.go) examples.

## Licensing

This project is governed by [Mizumoto General Public License v1.2](License/Mizumoto%20General%20Public%20License%20v1.2.md). Basically a Mozilla 2.0 public license, but with extra restrictions:

By using any part of this project, you are deemed to have fully understanding and acceptance of the following terms：

1. You must conspicuously display, without modification, this License and the notice on each redistributed or derivative copy of the License Covered Work.
2. Any non-independent developers companies/groups/legal entities or other organizations should ensure that employees are not oppressed or exploited, and that employees can always receive a reasonable salary for their legal working hours.
3. Any independent or non-independent developers/companies/groups/legal entities or other organizations, shall ensure that it has a clear conscience, including and not limited to **opposition to any form of Nazi or Neo-Nazism organization(s)**.

Otherwise these Individuals / Companies / Groups / Legal-entities **will not have the right to copy / modify / redistribute any code / file / algorithm** governed by MGPL v1.2.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	asciiBanner = `
	 ________      __________        .__                                    
	/  _____/  ____\______   \_____  |  | _____    ____   ____  ___________ 
   /   \  ___ /  _ \|    |  _/\__  \ |  | \__  \  /    \_/ ___\/ __ \_  __ \
   \    \_\  (  <_> )    |   \ / __ \|  |__/ __ \|   |  \  \__\  ___/|  | \/
	\______  /\____/|______  /(____  /____(____  /___|  /\___  >___  >__|   
		   \/              \/      \/          \/     \/     \/    \/       
`
)

type Config struct {
	Schema              string      `json:"schema"`
	SSLCertKey          string      `json:"ssl_cert_key"`
	SSLCert             string      `json:"ssl_cert"`
	Location            []*Location `json:"location"` // location of reverse proxy
	Port                int         `json:"port"`
	HealthCheck         bool        `json:"health_check"`
	HealthCheckInterval uint64      `json:"health_check_interval"`
	MaxConnections      uint64      `json:"max_connections"`
}

// Location routing information.
type Location struct {
	Pattern     string   `json:"pattern"`      // proxy route matching pattern
	ProxyPass   []string `json:"proxy_pass"`   // URL(s) of the reverse proxy
	BalanceMode string   `json:"balance_mode"` // load balancing algorithm
}

// LoadConfig loads the configuration from the given file.
func LoadConfig(file string) (*Config, error) {
	var config Config
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// PrintBanner prints the ASCII banner.
func PrintBanner() {
	fmt.Println(asciiBanner)
}

// PrintConfig prints the configuration.
func PrintConfig(config *Config) {
	fmt.Println("Schema:", config.Schema)
	fmt.Println("Port:", config.Port)
	fmt.Println("HealthCheck:", config.HealthCheck)
	fmt.Println("MaxConnections:", config.MaxConnections)
	PrintBanner()
	for _, location := range config.Location {
		fmt.Println("Pattern:", location.Pattern)
		fmt.Println("ProxyPass:", location.ProxyPass)
		fmt.Println("BalanceMode:", location.BalanceMode)
	}
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Schema != "http" && c.Schema != "https" {
		return fmt.Errorf("schema must be http or https")
	}
	if len(c.Location) == 0 {
		return fmt.Errorf("location must be specified")
	}
	if c.Schema == "https" && (len(c.SSLCert) == 0 || len(c.SSLCertKey) == 0) {
		return fmt.Errorf("ssl_cert and ssl_cert_key must be specified")
	}

	if c.HealthCheckInterval == 0 {
		return fmt.Errorf("health_check_interval must be greater than 0")
	}
	for _, location := range c.Location {
		if len(location.Pattern) == 0 {
			return fmt.Errorf("pattern must be specified")
		}
		if len(location.ProxyPass) == 0 {
			return fmt.Errorf("proxy_pass must be specified")
		}
		if len(location.BalanceMode) == 0 {
			return fmt.Errorf("balance_mode must be specified")
		}
	}
	return nil
}

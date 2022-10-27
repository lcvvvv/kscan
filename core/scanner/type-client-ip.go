package scanner

import (
	"kscan/lib/osping"
	"kscan/lib/tcpping"
	"net"
)

type IPClient struct {
	*client
	HandlerAlive func(addr net.IP)
	HandlerDie   func(addr net.IP)
	HandlerError func(addr net.IP, err error)
}

func NewIPScanner(config *Config) *IPClient {
	var client = &IPClient{
		client:       newConfig(config, config.Threads),
		HandlerAlive: func(addr net.IP) {},
		HandlerDie:   func(addr net.IP) {},
		HandlerError: func(addr net.IP, err error) {},
	}
	client.pool.Interval = config.Interval
	client.pool.Function = func(in interface{}) {
		ip := in.(net.IP)
		if client.config.HostDiscoverClosed == true {
			client.HandlerAlive(ip)
			return
		}
		if osping.Ping(ip.String()) == true {
			client.HandlerAlive(ip)
			return
		}
		if err := tcpping.PingPorts(ip.String(), config.Timeout); err == nil {
			client.HandlerAlive(ip)
			return
		}
		client.HandlerDie(ip)
	}
	return client
}

func (c *IPClient) Push(ips ...net.IP) {
	for _, ip := range ips {
		c.pool.Push(ip)
	}
}

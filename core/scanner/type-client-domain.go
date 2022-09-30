package scanner

import (
	"kscan/core/cdn"
	"net"
)

var CDNCheck = false

type DomainClient struct {
	*client
	HandlerIsCDN  func(domain, CDNInfo string)
	HandlerRealIP func(domain string, ip net.IP)
	HandlerError  func(domain string, err error)
}

func NewDomainScanner(config *Config) *DomainClient {
	var client = &DomainClient{
		client:        newConfig(config, config.Threads),
		HandlerIsCDN:  func(domain, CDNInfo string) {},
		HandlerRealIP: func(domain string, ip net.IP) {},
		HandlerError:  func(domain string, err error) {},
	}
	client.pool.Interval = config.Interval
	client.pool.Function = func(in interface{}) {
		domain := in.(string)

		ip, err := cdn.Resolution(domain)
		if err != nil {
			client.HandlerError(domain, err)
			return
		}
		if CDNCheck == false {
			client.HandlerRealIP(domain, net.ParseIP(ip))
			return
		}

		if ok, result, _ := cdn.FindWithDomain(domain); ok {
			client.HandlerIsCDN(domain, result)
			return
		}
		if ok, result, _ := cdn.FindWithIP(ip); ok {
			client.HandlerIsCDN(domain, result)
			return
		}
		client.HandlerRealIP(domain, net.ParseIP(ip))
	}
	return client
}

func (c *DomainClient) Push(domain string) {
	c.pool.In <- domain
}

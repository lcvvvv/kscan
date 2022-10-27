package scanner

import (
	"github.com/lcvvvv/gonmap"
	"net"
)

type foo1 struct {
	addr net.IP
	num  int
}

type PortClient struct {
	*client

	HandlerClosed     func(addr net.IP, port int)
	HandlerOpen       func(addr net.IP, port int)
	HandlerNotMatched func(addr net.IP, port int, response string)
	HandlerMatched    func(addr net.IP, port int, response *gonmap.Response)
	HandlerError      func(addr net.IP, port int, err error)
}

func NewPortScanner(config *Config) *PortClient {
	var client = &PortClient{
		client:            newConfig(config, config.Threads),
		HandlerClosed:     func(addr net.IP, port int) {},
		HandlerOpen:       func(addr net.IP, port int) {},
		HandlerNotMatched: func(addr net.IP, port int, response string) {},
		HandlerMatched:    func(addr net.IP, port int, response *gonmap.Response) {},
		HandlerError:      func(addr net.IP, port int, err error) {},
	}
	client.pool.Interval = config.Interval
	client.pool.Function = func(in interface{}) {
		nmap := gonmap.New()
		nmap.SetTimeout(config.Timeout)
		if config.DeepInspection == true {
			nmap.OpenDeepIdentify()
		}
		value := in.(foo1)
		status, response := nmap.ScanTimeout(value.addr.String(), value.num, 100*config.Timeout)
		switch status {
		case gonmap.Closed:
			client.HandlerClosed(value.addr, value.num)
		case gonmap.Open:
			client.HandlerOpen(value.addr, value.num)
		case gonmap.NotMatched:
			client.HandlerNotMatched(value.addr, value.num, response.Raw)
		case gonmap.Matched:
			client.HandlerMatched(value.addr, value.num, response)
		}
	}
	return client
}

func (c *PortClient) Push(ip net.IP, num int) {
	c.pool.Push(foo1{ip, num})
}

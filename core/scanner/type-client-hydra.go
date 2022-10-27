package scanner

import (
	"kscan/app"
	"kscan/core/hydra"
	"net"
)

type HydraClient struct {
	*client
	HandlerSuccess func(addr net.IP, port int, protocol string, auth *hydra.Auth)
	HandlerError   func(addr net.IP, port int, protocol string, err error)
}

type foo3 struct {
	ipAddr   net.IP
	port     int
	protocol string
}

func NewHydraScanner(config *Config) *HydraClient {
	var client = &HydraClient{
		client:         newConfig(config, config.Threads),
		HandlerSuccess: func(addr net.IP, port int, protocol string, auth *hydra.Auth) {},
		HandlerError:   func(addr net.IP, port int, protocol string, err error) {},
	}
	client.pool.Interval = config.Interval
	client.pool.Function = func(in interface{}) {
		value := in.(foo3)
		ipAddr := value.ipAddr
		protocol := value.protocol
		port := value.port
		//适配爆破模块
		authInfo := hydra.NewAuthInfo(ipAddr.String(), port, protocol)
		crack := hydra.NewCracker(authInfo, app.Setting.HydraUpdate, 10)
		auth, err := crack.Run()
		if err != nil {
			client.HandlerError(ipAddr, port, protocol, err)
		} else {
			client.HandlerSuccess(ipAddr, port, protocol, auth)
		}
	}
	return client
}

func (c *HydraClient) Push(addr net.IP, port int, protocol string) {
	c.pool.Push(foo3{addr, port, protocol})
}

package app

import (
	"time"
)

var Config = config{
	HostTarget: []string{},
	UrlTarget:  []string{},
	Path:       "/",
	Port:       WOOYUN_PORT_TOP_1000[:400],
	HttpCode:   []int{},
	Output:     nil,
	Proxy:      "",
	Host:       "",
	Threads:    500,
	Timeout:    time.Second * 5,
}

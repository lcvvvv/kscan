package app

var Config = config{
	HostTarget:    []string{},
	HostTargetNum: 0,
	UrlTarget:     []string{},
	UrlTargetNum:  0,
	PingAliveMap:  nil,
	Path:          "/",
	Port:          WOOYUN_PORT_TOP_1000[:400],
	PortNum:       0,
	Output:        nil,
	Proxy:         "",
	Host:          "",
	Threads:       500,
	Timeout:       0,
}

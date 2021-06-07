package app

import "runtime"

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
	Encoding:      "utf-8",
	OSEncoding:    getOSEncoding(),
	NewLine:       getNewline(),
}

func getNewline() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	} else {
		return "\n"
	}
}

func getOSEncoding() string {
	if runtime.GOOS == "windows" {
		return "gb2312"
	} else {
		return "utf-8"
	}
}

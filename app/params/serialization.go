package params

import (
	"../../app/config"
	"../../lib/IP"
	"../../lib/misc"
	"../../lib/urlparse"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type serParams struct {
	HostTarget, UrlTarget []string
	Port, HttpCode        []int
	Path                  []string
	Output, Proxy, Host   string
	Threads               int
}

var SerParams = serParams{
	HostTarget: nil,
	UrlTarget:  nil,
	Port:       nil,
	HttpCode:   nil,
	Path:       nil,
	Output:     "",
	Proxy:      "",
	Host:       "",
	Threads:    0,
}

func serializationParams() {
	serializationParamsTarget(params.target)
	serializationParamsPort()
	serializationParamsHttpCode()
	serializationParamsPath()
	serializationParamsOutput()
	serializationParamsProxy()
	serializationParamsHosts()
	serializationParamsThreads()
	//fmt.Print(SerParams)
}

func serializationParamsTarget(t string) {
	var HostTarget = SerParams.HostTarget
	var UrlTarget = SerParams.UrlTarget
	isFile, _ := regexp.MatchString("^file:.+$", t)
	if isFile {
		t = strings.Replace(t, "file:", "", 1)
		err := misc.ReadLine(t, serializationParamsTarget)
		if err != nil {
			fmt.Print(err)
			os.Exit(0)
		}
		return
	}
	//判断target字符串是否为类IP/MASK
	Hosts, err := IP.IPMask2IPArr(t)
	if err == nil {
		SerParams.HostTarget = append(HostTarget, Hosts...)
		return
	}
	//判断target字符串是否为类URL
	v, err := urlparse.Load(t)
	if err == nil {
		if v.Scheme != "" {
			SerParams.UrlTarget = misc.UniStrAppend(UrlTarget, t)
			SerParams.HostTarget = misc.UniStrAppend(HostTarget, v.Host)
			return
		} else {
			SerParams.HostTarget = misc.UniStrAppend(HostTarget, v.Host)
			return
		}
	}

}

func serializationParamsPort() {
	if params.port != "" {
		strPorts := strings.Split(params.port, ",")
		serStrPorts, _ := misc.StrArr2IntArr(strPorts)
		SerParams.Port = serStrPorts
		return
	}
	if params.top != 0 {
		strPorts := config.Config.Ports[0:params.top]
		SerParams.Port = strPorts
		return
	}
	if params.port == "" && params.top == 0 {
		strPorts := config.Config.Ports[0:config.Config.Top]
		SerParams.Port = strPorts
		return
	}
}

func serializationParamsHttpCode() {
	strHttpCodes := strings.Split(params.httpCode, ",")
	serStrHttpCodes, _ := misc.StrArr2IntArr(strHttpCodes)
	SerParams.HttpCode = serStrHttpCodes
}

func serializationParamsPath() {
	var fixStrPaths []string
	strPaths := strings.Split(params.path, ",")
	for _, path := range strPaths {
		path = fixPath(path)
		fixStrPaths = append(fixStrPaths, path)
	}
	SerParams.Path = fixStrPaths
}

func fixPath(path string) string {
	if path == "" {
		return path
	}
	if path[0:1] == "/" {
		return fixPath(path[1:])
	} else {
		return path
	}
}

func serializationParamsOutput() {
	SerParams.Output = params.output
}

func serializationParamsProxy() {
	SerParams.Proxy = params.proxy

}

func serializationParamsHosts() {
	SerParams.Host = params.host
}

func serializationParamsThreads() {
	SerParams.Threads = params.threads
}

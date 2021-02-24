package params

import (
	"app/config"
	"lib/IP"
	"lib/misc"
	"lib/slog"
	"lib/urlparse"
	"os"
	"regexp"
	"strings"
)

type serParams struct {
	HostTarget, UrlTarget []string
	Port, HttpCode        []int
	Path                  []string
	Output, Proxy, Host   string
	Threads, Timeout      int
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
	Timeout:    3,
}

func serializationParams() {
	serializationParamsTarget(Params.target)
	serializationParamsPort()
	serializationParamsHttpCode()
	serializationParamsPath()
	serializationParamsOutput()
	serializationParamsProxy()
	serializationParamsHosts()
	serializationParamsThreads()
	serializationParamsTimeout()
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
			slog.Error(err.Error())
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
	if Params.port != "" {
		SerParams.Port = intParam2IntArr(Params.port)
		return
	}
	if Params.top != 0 {
		strPorts := config.Config.Ports[0:Params.top]
		SerParams.Port = strPorts
		return
	}
	if Params.port == "" && Params.top == 0 {
		strPorts := config.Config.Ports[0:config.Config.Top]
		SerParams.Port = strPorts
		return
	}
}

func serializationParamsHttpCode() {
	SerParams.HttpCode = intParam2IntArr(Params.httpCode)
}

func serializationParamsPath() {
	var fixStrPaths []string
	strPaths := strings.Split(Params.path, ",")
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
	SerParams.Output = Params.output
}

func serializationParamsProxy() {
	SerParams.Proxy = Params.proxy

}

func serializationParamsHosts() {
	SerParams.Host = Params.host
}

func serializationParamsThreads() {
	SerParams.Threads = Params.threads
}

func serializationParamsTimeout() {
	SerParams.Timeout = Params.timeout
}

func intParam2IntArr(v string) []int {
	var res []int
	vArr := strings.Split(v, ",")
	for _, v := range vArr {
		var vvArr []int
		if strings.Contains(v, "-") {
			iArr := strings.Split(v, "-")
			if len(iArr) != 2 {
				slog.Error("参数输入错误！！！")
				os.Exit(0)
			} else {
				smallNum := misc.Str2Int(iArr[0])
				bigNum := misc.Str2Int(iArr[1])
				if smallNum >= bigNum {
					slog.Error("参数输入错误！！！")
					os.Exit(0)
				}
				vvArr = append(vvArr, makeIntList(smallNum, bigNum)...)
			}
		} else {
			vvArr = append(vvArr, misc.Str2Int(v))
		}
		res = append(res, vvArr...)
	}
	return res
}

func makeIntList(s int, b int) []int {
	var iArr []int
	for i := 0; i <= b-s; i++ {
		iArr = append(iArr, i+s)
	}
	return iArr
}

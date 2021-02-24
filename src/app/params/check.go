package params

import (
	"lib/misc"
	"lib/slog"
	"os"
	"regexp"
)

var OutPutFile *os.File

func initOutPutFile(path string) *os.File {
	if misc.FileIsExist(path) {
		f, err := os.Open(path)
		if err != nil {
			slog.Error(err.Error())
		}
		return f
	} else {
		f, err := os.Create(path)
		if err != nil {
			slog.Error(err.Error())
		}
		return f
	}
}

func checkParams() {
	//判断冲突参数
	if params.port != "" && params.top != 0 {
		slog.Error("PORT、TOP只允许同时出现一个")
	}
	//判断内容
	if params.target != "" {
		if params.port != "" {
			if !checkIntsParam(params.port) {
				slog.Error("PORT参数输入错误,其格式应为80，8080，8081-8090")
				os.Exit(0)
			}
		}
		if params.top != 0 {
			if params.top > 1000 || params.top < 1 {
				slog.Error("TOP参数输入错误,TOP参数应为1-1000之间的整数。")
				os.Exit(0)
			}
		}
		if params.output != "" {
			//验证output参数
			f := initOutPutFile(params.output)
			OutPutFile = f
		}
		if params.proxy != "" {
			if !checkProxyParam(params.proxy) {
				slog.Error("PROXY参数输入错误，其格式应为：http://IP:PORT，支持socks5/4")
				os.Exit(0)
			}
		}
		if params.path != "" {
			if !checkStringsParam(params.path) {
				slog.Error("PATH参数输入错误，其格式应为：/asdfasdf，可使用逗号输入多个路径")
				os.Exit(0)
			}
		}
		if params.host != "" {
			//验证host参数
		}
		if params.threads != 0 {
			//验证threads参数
		}
		if params.timeout != 3 {
			//验证timeout参数
		}
		if params.httpCode != "" {
			if !checkIntsParam(params.httpCode) {
				slog.Error("HTTPCODE参数输入错误，其格式应为200可用逗号输入多个状态码")
				os.Exit(0)
			}
		}
	} else {
		slog.Error("必须输入TARGET参数")
		os.Exit(0)
	}
}

func checkIntsParam(v string) bool {
	matched, err := regexp.MatchString("^((?:[0-9])+(?:-[0-9]+)?)(?:,(?:[0-9])+-(?:[0-9])+)*$", v)
	if err != nil {
		return false
	}
	if matched {
		return true
	} else {
		return false
	}
}

func checkStringsParam(v string) bool {
	matched, err := regexp.MatchString("^([A-Za-z0-9/]+)(,[A-Za-z0-9/])*$", v)
	if err != nil {
		return false
	}
	if matched {
		return true
	} else {
		return false
	}
}

func checkProxyParam(v string) bool {
	matched, err := regexp.MatchString("^(http|https|socks5|socks4)://[0-9.]+:[0-9]+$", v)
	if err != nil {
		return false
	}
	if matched {
		return true
	} else {
		return false
	}
}

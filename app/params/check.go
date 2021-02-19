package params

import (
	"../../lib/misc"
	"fmt"
	"os"
	"regexp"
)

var OutPutFile *os.File

func initOutPutFile(path string) *os.File {
	if misc.FileIsExist(path) {
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("[-]%s", err.Error())
			os.Exit(0)
		}
		return f
	} else {
		f, err := os.Create(path)
		if err != nil {
			fmt.Printf("[-]%s", err.Error())
			os.Exit(0)
		}
		return f
	}
}

func checkParams() {
	//判断冲突参数
	if params.port != "" && params.top != 0 {
		fmt.Print("port、top只允许同时出现一个")
		os.Exit(0)
	}
	//判断内容
	if params.target != "" {
		if params.port != "" {
			if !checkIntsParam(params.port) {
				fmt.Print("port参数输入错误")
				os.Exit(0)
			}
		}
		if params.top != 0 {
			if params.top > 1000 || params.top < 1 {
				fmt.Print("top参数输入错误")
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
				fmt.Print("proxy参数输入错误")
				os.Exit(0)
			}
		}
		if params.path != "" {
			if !checkStringsParam(params.path) {
				fmt.Print("path参数输入错误")
				os.Exit(0)
			}
		}
		if params.host != "" {
			//验证host参数
		}
		if params.threads != 0 {
			//验证threads参数
		}
		if params.threads != 3 {
			//验证threads参数
		}
		if params.httpCode != "" {
			if !checkIntsParam(params.httpCode) {
				fmt.Print("port参数输入错误")
				os.Exit(0)
			}
		}
	} else {
		fmt.Print("必须输入target参数")
		os.Exit(0)
	}
}

func checkIntsParam(v string) bool {
	matched, err := regexp.MatchString("^([0-9]+,*)+[0-9]$", v)
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

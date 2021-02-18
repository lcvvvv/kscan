package params

import (
	"../config"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Params struct {
	help                                              bool
	target, port, output, proxy, path, host, httpCode string
	top, threads                                      int
}

var params Params

//logo信息
const logo = `
 _  __ _____  _____     *     _   _
|#|/#//####/ /#####|   /#\   |#\ |#|
|#.#/|#|___  |#|      /###\  |##\|#|
|##|  \#####\|#|     /#/_\#\ |#.#.#|
|#.#\_____|#||#|____/#/###\#\|#|\##|
|#|\#\#####/ \#####/#/ v0.1\#\#| \#|
轻量资产测绘工具                by：kv2
`

//帮助信息
const help = `
optional arguments:
  -h , --help     show this help message and exit
  -t , --target   直接扫描指定对象,支持IP、URL、IP/[16-32]、file:/tmp/target.txt
  -p , --port     扫描指定端口，默认会扫描
  -o , --output   将扫描结果保存到文件
  --top           扫描WooYun统计开放端口前x个，最高支持1000个
  --proxy         设置代理{socks5/socks4/https/http}://IP:port
  --threads       线程参数
  --http-code     指定会记录的HTTP状态码，逗号分割,默认会记录200,301,302,403,404
  --path          指定请求访问的目录，逗号分割，慎用！
  --host          指定所有请求的头部HOST值，慎用！
`

const usage = "usage: kscan [-h,--help] (-t,--target) [-p,--port|--top] [-o,--output] [--proxy] [--threads] [--http-code] [--path] [--host]\n"

//初始化参数
func initParams() {
	//自定义Usage
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, logo)
	}
	flag.BoolVar(&params.help, "h", false, "")
	flag.BoolVar(&params.help, "help", false, "")
	flag.StringVar(&params.target, "t", "", "")
	flag.StringVar(&params.target, "target", "", "")
	flag.StringVar(&params.port, "p", "", "")
	flag.StringVar(&params.port, "port", "", "")
	flag.StringVar(&params.output, "o", "", "")
	flag.StringVar(&params.output, "output", "", "")
	flag.StringVar(&params.proxy, "proxy", "", "")
	flag.StringVar(&params.path, "path", "", "")
	flag.StringVar(&params.host, "host", "", "")
	flag.StringVar(&params.httpCode, "http-code", "", "")
	flag.IntVar(&params.top, "top", 0, "")
	flag.IntVar(&params.threads, "threads", 0, "")
}

func LoadParams() {
	initParams()
	flag.Parse()
	//不带参数则对应usage
	if len(os.Args) == 1 {
		_, _ = fmt.Fprintf(os.Stderr, logo)
		_, _ = fmt.Fprintf(os.Stderr, usage)
		os.Exit(0)
	}
	if params.help {
		_, _ = fmt.Fprintf(os.Stderr, logo)
		_, _ = fmt.Fprintf(os.Stderr, usage)
		_, _ = fmt.Fprintf(os.Stderr, help)
		os.Exit(0)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, logo)
		checkParams()
		//加载配置文件
		config.LoadConfig()
		if params.top == 0 {
			params.top = config.Config.Top
		}
		if params.threads == 0 {
			params.threads = config.Config.Threads
		}
		if params.path == "" {
			params.path = strings.Join(config.Config.Path, ",")
		}
		serializationParams()
	}
}

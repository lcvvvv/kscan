package params

import (
	"flag"
	"fmt"
	"kscan/src/app/config"
	"os"
)

type OsArgs struct {
	help, Debug                                       bool
	target, port, output, proxy, path, host, httpCode string
	top, threads, timeout                             int
}

var Params OsArgs

//logo信息
const logo = `
 _  __ _____  _____     *     _   _
|#|/#//####/ /#####|   /#\   |#\ |#|
|#.#/|#|___  |#|      /###\  |##\|#|
|##|  \#####\|#|     /#/_\#\ |#.#.#|
|#.#\_____|#||#|____/#/###\#\|#|\##|
|#|\#\#####/ \#####/#/ v1.12#\#| \#|
           轻量级资产测绘工具 by：kv2

`

//帮助信息
const help = `
optional arguments:
  -h , --help     show this help message and exit
  -t , --target   直接扫描指定对象,支持IP、URL、IP/[16-32]、file:/tmp/target.txt
  -p , --port     扫描指定端口，默认会扫描TOP400，支持：80,8080,8088-8090
  -o , --output   将扫描结果保存到文件
  --top           扫描WooYun统计开放端口前x个，最高支持1000个
  --proxy         设置代理(socks5|socks4|https|http)://IP:Port
  --threads       线程参数,默认线程4000
  --http-code     指定会记录的HTTP状态码，逗号分割,默认会记录200,301,302,403,404
  --path          指定请求访问的目录，逗号分割，慎用！
  --host          指定所有请求的头部HOSTS值，慎用！
  --timeout       设置超时时间，默认3秒钟，单位为秒！

`

const usage = "usage: kscan [-h,--help] (-t,--target) [-p,--port|--top] [-o,--output] [--proxy] [--threads] [--http-code] [--path] [--host] [--timeout]\n\n"

//初始化函数
func Init() {
	initParams()
	flag.Parse()
	//不带参数则对应usage
	if len(os.Args) == 1 {
		fmt.Print(logo)
		fmt.Print(usage)
		os.Exit(0)
	}
	if Params.help {
		fmt.Print(logo)
		fmt.Print(usage)
		fmt.Print(help)
		os.Exit(0)
	}
	//打印logo
	fmt.Print(logo)
	checkParams()
}

//初始化参数
func initParams() {
	//自定义Usage
	flag.Usage = func() {
		fmt.Print(logo)
	}
	flag.BoolVar(&Params.help, "h", false, "")
	flag.BoolVar(&Params.help, "help", false, "")
	flag.BoolVar(&Params.Debug, "debug", false, "")
	flag.BoolVar(&Params.Debug, "d", false, "")
	flag.StringVar(&Params.target, "t", "", "")
	flag.StringVar(&Params.target, "target", "", "")
	flag.StringVar(&Params.port, "p", "", "")
	flag.StringVar(&Params.port, "port", "", "")
	flag.StringVar(&Params.output, "o", "", "")
	flag.StringVar(&Params.output, "output", "", "")
	flag.StringVar(&Params.proxy, "proxy", "", "")
	flag.StringVar(&Params.path, "path", "", "")
	flag.StringVar(&Params.host, "host", "", "")
	flag.StringVar(&Params.httpCode, "http-code", "", "")
	flag.IntVar(&Params.top, "top", 400, "")
	flag.IntVar(&Params.threads, "threads", 400, "")
	flag.IntVar(&Params.timeout, "timeout", 3, "")
}

func (o OsArgs) Target() string {
	return o.target
}
func (o OsArgs) Port() string {
	return o.port
}
func (o OsArgs) Output() string {
	return o.output
}
func (o OsArgs) Proxy() string {
	return o.proxy
}
func (o OsArgs) Path() string {
	return o.path
}
func (o OsArgs) Host() string {
	return o.host
}
func (o OsArgs) HttpCode() string {
	return o.httpCode
}
func (o OsArgs) Top() int {
	return o.top
}
func (o OsArgs) Threads() int {
	return o.threads
}
func (o OsArgs) Timeout() int {
	return o.timeout
}

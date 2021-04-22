package params

import (
	"flag"
	"fmt"
	"os"
)

type OsArgs struct {
	help, Debug, scanPing                   bool
	target, port, output, proxy, path, host string
	top, threads, timeout                   int
}

var Params OsArgs

//logo信息
const logo = `
 _  __ _____  _____     *     _   _
|#|/#//####/ /#####|   /#\   |#\ |#|
|#.#/|#|___  |#|      /###\  |##\|#|
|##|  \#####\|#|     /#/_\#\ |#.#.#|
|#.#\_____|#||#|____/#/###\#\|#|\##|
|#|\#\#####/ \#####/#/ v1.14#\#| \#|
           轻量级资产测绘工具 by：kv2

`

//帮助信息
const help = `
optional arguments:
  -h , --help     show this help message and exit
  --ping          在扫描端口之前会先进行Ping探测，若不存活，则不会进行端口扫描
  -t , --target   指定探测对象：
                  IP地址：114.114.114.114
                  IP地址段：114.114.114.114/24,不建议子网掩码小于12
                  URL地址：https://www.baidu.com
                  文件地址：file:/tmp/target.txt
  -p , --port     扫描指定端口，默认会扫描TOP400，支持：80,8080,8088-8090
  -o , --output   将扫描结果保存到文件
  --top           扫描WooYun统计开放端口前x个，最高支持1000个
  --proxy         设置代理(socks5|socks4|https|http)://IP:Port
  --threads       线程参数,默认线程400,最大值为2048
  --path          指定请求访问的目录，逗号分割，慎用！
  --host          指定所有请求的头部HOSTS值，慎用！
  --timeout       设置超时时间，默认为预设的探针超时时间！
`

const usage = "usage: kscan [-h,--help] (-t,--target) [-p,--port|--top] [-o,--output] [--proxy] [--threads] [--path] [--host] [--timeout] [--ping]\n\n"

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
	flag.BoolVar(&Params.scanPing, "ping", false, "")
	flag.StringVar(&Params.target, "t", "", "")
	flag.StringVar(&Params.target, "target", "", "")
	flag.StringVar(&Params.port, "p", "", "")
	flag.StringVar(&Params.port, "port", "", "")
	flag.StringVar(&Params.output, "o", "", "")
	flag.StringVar(&Params.output, "output", "", "")
	flag.StringVar(&Params.proxy, "proxy", "", "")
	flag.StringVar(&Params.path, "path", "", "")
	flag.StringVar(&Params.host, "host", "", "")
	flag.IntVar(&Params.top, "top", 400, "")
	flag.IntVar(&Params.threads, "threads", 400, "")
	flag.IntVar(&Params.timeout, "timeout", 0, "")
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
func (o OsArgs) Top() int {
	return o.top
}
func (o OsArgs) Threads() int {
	return o.threads
}
func (o OsArgs) Timeout() int {
	return o.timeout
}
func (o OsArgs) ScanPing() bool {
	return o.scanPing
}

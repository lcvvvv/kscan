package main

import (
	"kscan/app"
	"kscan/lib/gonmap"
	"kscan/lib/httpfinger"
	"kscan/lib/params"
	"kscan/lib/slog"
	"kscan/lib/spy"
	"kscan/run"
	"runtime"
	"time"
)

//logo信息
const logo = `
 _  __ _____  _____     *     _   _
|#|/#//####/ /#####|   /#\   |#\ |#|
|#.#/|#|___  |#|      /###\  |##\|#|
|##|  \#####\|#|     /#/_\#\ |#.#.#|
|#.#\_____|#||#|____/#/###\#\|#|\##|
|#|\#\#####/ \#####/#/ v1.26#\#| \#|
           轻量级资产测绘工具 by：kv2

`

//帮助信息
const help = `
optional arguments:
  -h , --help     show this help message and exit
  -t , --target   指定探测对象：
                  IP地址：114.114.114.114
                  IP地址段：114.114.114.114/24,不建议子网掩码小于12
                  IP地址段：114.114.114.114-115.115.115.115
                  URL地址：https://www.baidu.com
                  文件地址：file:/tmp/target.txt
  -p , --port     扫描指定端口，默认会扫描TOP400，支持：80,8080,8088-8090
  -o , --output   将扫描结果保存到文件
  -Pn          	  使用此参数后，将不会进行智能存活性探测，现在默认会开启智能存活性探测，提高效率
  --check         针对目标地址做指纹识别，仅不会进行端口探测
  --top           扫描经过筛选处理的常见端口TopX，最高支持1000个，默认为TOP4000
  --proxy         设置代理(socks5|socks4|https|http)://IP:Port
  --threads       线程参数,默认线程400,最大值为2048
  --path          指定请求访问的目录，逗号分割
  --host          指定所有请求的头部Host值
  --timeout       设置超时时间
  --encoding      设置终端输出编码，可指定为：gb2312、utf-8
  --spy           网段探测模式，此模式下将自动探测主机可达的内网网段,无需配置其他任何参数
  --rarity        指定Nmap指纹识别级别[0-9],数字越大可识别的协议越多越准确，但是扫描时间会更长,默认为：4
`

const usage = "usage: kscan [-h,--help] (-t,--target) [--spy] [-p,--port|--top] [-o,--output] [--proxy] [--threads] [--path] [--host] [--timeout] [-Pn] [--check] [--encoding] [--rarity]\n\n"

func main() {
	startTime := time.Now()

	//环境初始化
	Init()

	//校验升级情况
	//app.CheckUpdate()
	if app.Setting.Spy {
		spy.Start()
	} else {
		//扫描模块初始化
		KscanInit()
		//开始扫描
		run.Start(app.Setting)
	}
	//计算程序运行时间
	elapsed := time.Since(startTime)
	slog.Infof("程序执行总时长为：[%s]", elapsed.String())
	slog.Info("若有问题欢迎来我的Github提交Bug[https://github.com/lcvvvv/kscan/]")
}

func Init() {
	param := params.New(logo, usage, help)
	//参数初始化
	param.LoadOsArgs()
	//日志初始化
	slog.Init(param.Debug(), param.Encoding())
	//输出Banner
	param.PrintBanner()
	//参数合法性校验
	param.CheckArgs()
	//配置文件初始化
	app.Setting.Load(param)
	slog.Warning("当前环境为：", runtime.GOOS, ", 输出编码为：", app.Setting.Encoding)
}

func KscanInit() {
	slog.Warning("开始读取扫描对象...")
	slog.Infof("成功读取URL地址:[%d]个\n", len(app.Setting.UrlTarget))

	if app.Setting.Check == false {
		slog.Infof("成功读取主机地址:[%d]个，待检测端口:[%d]个\n", len(app.Setting.HostTarget), len(app.Setting.HostTarget)*len(app.Setting.Port))
	}
	//HTTP指纹库初始化
	r := httpfinger.Init()
	slog.Infof("成功加载favicon指纹:[%d]条，keyword指纹:[%d]条\n", r["FaviconHash"], r["KeywordFinger"])
	//gonmap探针/指纹库初始化
	r = gonmap.Init(app.Setting.Rarity, app.Setting.Timeout)
	slog.Infof("成功加载NMAP探针:[%d]个,指纹[%d]条\n", r["PROBE"], r["MATCH"])
	slog.Warningf("本次扫描将使用NMAP探针:[%d]个,指纹[%d]条\n", r["USED_PROBE"], r["USED_MATCH"])
}

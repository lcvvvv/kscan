package main

import (
	"github.com/lcvvvv/gonmap"
	"kscan/app"
	"kscan/lib/httpfinger"
	"kscan/lib/params"
	"kscan/lib/slog"
	"kscan/run"
	"time"
)

func main() {
	startTime := time.Now()
	//参数初始化
	params.Init()
	//日志初始化
	slog.Init(params.Params.Debug)
	//参数合法性校验
	params.CheckParams()
	//配置文件初始化
	app.Config.Load(params.Params)
	slog.Warning("开始读取扫描对象...")
	slog.Infof("成功读取URL地址:[%d]个\n", len(app.Config.UrlTarget))
	slog.Infof("成功读取主机地址:[%d]个，待检测端口:[%d]个\n", len(app.Config.HostTarget), len(app.Config.HostTarget)*len(app.Config.Port))
	//HTTP指纹库初始化
	r := httpfinger.Init()
	slog.Infof("成功加载favicon指纹:[%d]条，keyword指纹:[%d]条\n", r["FaviconHash"], r["KeywordFinger"])
	//加载gonmap探针/指纹库
	r = gonmap.Init(5, app.Config.Timeout)
	slog.Infof("成功加载NMAP探针:[%d]个,指纹[%d]条\n", r["PROBE"], r["MATCH"])
	slog.Warningf("本次扫描将使用NMAP探针:[%d]个,指纹[%d]条\n", r["USED_PROBE"], r["USED_MATCH"])

	//校验升级情况
	//app.CheckUpdate()

	//开始扫描
	run.Start()
	//计算程序运行时间
	elapsed := time.Since(startTime)
	slog.Infof("程序执行总时长为：[%s]", elapsed.String())
}

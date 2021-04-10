package main

import (
	"github.com/lcvvvv/gonmap"
	"kscan/app"
	"kscan/lib/httpfinger"
	"kscan/lib/params"
	"kscan/lib/run"
	"kscan/lib/slog"
)

func main() {
	//参数初始化
	params.Init()
	//日志初始化
	slog.Init(params.Params.Debug)
	//配置文件初始化
	app.Config.Load(params.Params)
	slog.Warning("开始读取扫描对象...")
	slog.Infof("成功读取URL地址:[%d]个\n", len(app.Config.UrlTarget))
	slog.Infof("成功读取主机地址:[%d]个，待检测端口:[%d]个\n", len(app.Config.HostTarget), len(app.Config.HostTarget)*len(app.Config.Port))
	//指纹库初始化
	r := httpfinger.Init()
	slog.Infof("成功加载favicon指纹:[%d]条，keyword指纹:[%d]条\n", r["FaviconHash"], r["KeywordFinger"])
	//加载gonmap探针/指纹库
	r = gonmap.Init(5)
	slog.Infof("成功加载探针:[%d]个,指纹[%d]条\n", r["PROBE"], r["MATCH"])
	slog.Warningf("本次扫描将使用探针:[%d]个,指纹[%d]条\n", r["USED_PROBE"], r["USED_MATCH"])

	//校验升级情况
	//app.CheckUpdate()

	//初始化可访问URL地址队列
	slog.Warning("正在压入URL地址队列...")
	run.InitHostPortQueue()
	//初始化端口扫描队列
	slog.Warning("正在压入端口扫描队列...")
	run.InitPortQueue()
	//开始扫描所有开放端口
	slog.Warning("开始扫描所有开放端口...")
	//run.ScanOpenPort()
	//开始获取所有开放端口的Banner
	run.Start()
}

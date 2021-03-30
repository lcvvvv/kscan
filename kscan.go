package main

import (
	"kscan/src/app/finger"
	"kscan/src/app/params"
	"kscan/src/app/run"
	"kscan/src/app/update"
	"kscan/src/lib/slog"
)

func main() {
	//环境初始化
	initEnv()
	//加载程序运行参数
	params.LoadParams()
	//校验升级情况
	update.CheckUpdate()
	//加载指纹数据
	//var KeywordFingers,HashFingers = finger.LoadFinger()
	finger.LoadFinger()
	//初始化可访问URL地址队列
	slog.Warning("正在压入URL地址队列...")
	run.InitUrlQueue()
	//初始化端口扫描队列
	slog.Warning("正在压入端口扫描队列...")
	run.InitPortQueue()
	//开始扫描所有开放端口
	slog.Warning("开始扫描所有开放端口...")
	//run.ScanOpenPort()
	//开始获取所有开放端口的Banner
	run.GetBanner()
}

func initEnv() {
	//参数初始化
	params.Init()
	//日志初始化
	slog.Init(params.Params.Debug)
	//sysType := runtime.GOOS
	//if sysType == "linux" {
	//	// LINUX系统
	//}
	//
	//if sysType == "windows" {
	//	// windows系统
	//}
	//
	//if sysType == "darwin" {
	//	// MAC系统
	//	exec.Command("ulimit","-n","10240")
	//}
}

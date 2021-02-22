package main

import (
	"app/finger"
	"app/params"
	"app/run"
	"app/update"
	"fmt"
)

func main() {
	//初始化
	initEnv()
	//加载程序运行参数
	params.LoadParams()
	//校验升级情况
	update.CheckUpdate()
	//加载指纹数据
	//var KeywordFingers,HashFingers = finger.LoadFinger()
	finger.LoadFinger()
	//初始化可访问URL地址队列
	fmt.Print("[*]正在压入URL地址队列...\n")
	run.InitUrlQueue()
	//初始化端口扫描队列
	fmt.Print("[*]正在压入端口扫描队列...\n")
	run.InitPortQueue()
	//开始扫描所有开放端口
	fmt.Print("[*]开始扫描所有开放端口...\n")
	//run.ScanOpenPort()
	//开始获取所有开放端口的Banner
	run.GetBanner()
}

func initEnv() {
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

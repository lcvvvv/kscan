package main

import (
	"./app/finger"
	"./app/params"
	"./app/run"
	"fmt"
)

func main() {
	//加载程序运行参数
	params.LoadParams()
	//加载指纹数据
	//var KeywordFingers,HashFingers = finger.LoadFinger()
	finger.LoadFinger()
	//初始化端口扫描队列
	fmt.Print("[*]正在初始化端口扫描队列...\n")
	run.InitPortQueue()
	//初始化可访问URL地址队列
	//fmt.Print("[*]正在初始化可访问URL地址队列...\n")
	run.InitUrlQueue()
	//开始扫描所有开放端口
	fmt.Print("[*]开始扫描所有开放端口...\n")
	//run.ScanOpenPort()
	//开始获取所有开放端口的Banner
	run.GetBanner()
}

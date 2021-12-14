package app

import (
	"fmt"
	"os"

	"kscan/lib/sflag"
)

type args struct {
	USAGE, HELP, LOGO, SYNTAX string

	Help, Debug, ClosePing, Check, CloseColor, Scan   bool
	Target, Port, Output, Proxy, Path, Host, Encoding string
	OutputJson                                        string
	Spy, Touch                                        string
	Top, Threads, Timeout                             int
	//hydra模块
	Hydra, HydraUpdate             bool
	HydraUser, HydraPass, HydraMod string
	//fofa模块
	Fofa, FofaField, FofaFixKeyword string
	FofaSize                        int
	FofaSyntax                      bool
}

var Args = args{}

//初始化参数
func (o *args) Parse() {
	//自定义Usage
	sflag.SetUsage(o.LOGO)
	//定义参数
	o.define()
	//实例化参数值
	sflag.Parse()
	//输出LOGO
	o.printBanner()
	//参数值合法性校验
	o.checkArgs()
}

//定义参数
func (o *args) define() {
	sflag.BoolVar(&o.Help, "h", false)
	sflag.BoolVar(&o.Help, "help", false)
	sflag.BoolVar(&o.Debug, "debug", false)
	sflag.BoolVar(&o.Debug, "d", false)
	//spy模块
	sflag.AutoVarString(&o.Spy, "spy", "None")
	//touch模块
	sflag.StringVar(&o.Touch, "touch", "None")
	//hydra模块
	sflag.BoolVar(&o.Hydra, "hydra", false)
	sflag.BoolVar(&o.HydraUpdate, "hydra-update", false)
	sflag.StringVar(&o.HydraUser, "hydra-user", "")
	sflag.StringVar(&o.HydraPass, "hydra-pass", "")
	sflag.StringVar(&o.HydraMod, "hydra-mod", "")
	//fofa模块
	sflag.StringVar(&o.Fofa, "fofa", "")
	sflag.StringVar(&o.Fofa, "f", "")
	sflag.StringVar(&o.FofaField, "fofa-field", "")
	sflag.StringVar(&o.FofaFixKeyword, "fofa-fix-keyword", "")
	sflag.IntVar(&o.FofaSize, "fofa-size", 100)
	sflag.BoolVar(&o.FofaSyntax, "fofa-syntax", false)
	sflag.BoolVar(&o.Scan, "scan", false)
	//kscan模块
	sflag.StringVar(&o.Target, "target", "")
	sflag.StringVar(&o.Target, "t", "")
	sflag.StringVar(&o.Port, "p", "")
	sflag.StringVar(&o.Port, "port", "")
	sflag.StringVar(&o.Proxy, "proxy", "")
	sflag.StringVar(&o.Path, "path", "")
	sflag.StringVar(&o.Host, "host", "")
	sflag.IntVar(&o.Top, "top", 400)
	sflag.IntVar(&o.Threads, "threads", 100)
	sflag.IntVar(&o.Timeout, "timeout", 3)
	sflag.BoolVar(&o.ClosePing, "Pn", false)
	sflag.BoolVar(&o.Check, "check", false)
	//输出模块
	sflag.StringVar(&o.Encoding, "encoding", "utf-8")
	sflag.StringVar(&o.Output, "o", "")
	sflag.StringVar(&o.Output, "output", "")
	sflag.StringVar(&o.OutputJson, "oJ", "")
	sflag.BoolVar(&o.CloseColor, "Cn", false)
}

func (o *args) SetLogo(logo string) {
	o.LOGO = logo
}

func (o *args) SetUsage(usage string) {
	o.USAGE = usage
}

func (o *args) SetSyntax(syntax string) {
	o.SYNTAX = syntax
}

func (o *args) SetHelp(help string) {
	o.HELP = help
}

//校验参数真实性
func (o *args) checkArgs() {
	//判断必须的参数是否存在
	if o.Target == "" && o.Fofa == "" && o.Spy == "None" && o.Touch == "None" {
		fmt.Print("至少有target、fofa、spy、touch参数中的一个")
		os.Exit(0)
	}
	//判断冲突参数
	if o.Target != "" && o.Fofa != "" && o.Spy != "None" && o.Touch == "None" {
		fmt.Print("target、fofa、spy、touch不能同时使用")
		os.Exit(0)
	}
	if o.Port != "" && o.Top != 400 {
		fmt.Print("port、top参数不能同时使用")
		os.Exit(0)
	}
	//判断内容
	if o.Port != "" && sflag.MultipleIntVerification(o.Port) == false {
		fmt.Print("PORT参数输入错误,其格式应为80，8080，8081-8090")
		os.Exit(0)
	}
	if o.Top != 0 && (o.Top > 1000 || o.Top < 1) {
		fmt.Print("TOP参数输入错误,TOP参数应为1-1000之间的整数。")
		os.Exit(0)
	}
	if o.Proxy != "" && sflag.ProxyStrVerification(o.Proxy) {
		fmt.Print("PROXY参数输入错误，其格式应为：http://ip:port，支持socks5/4")
	}
	if o.Path != "" && sflag.MultipleStrVerification(o.Path) {
		fmt.Print("PATH参数输入错误，其格式应为：/asdfasdf，可使用逗号输入多个路径")
	}
	if o.Threads != 0 && o.Threads > 2048 {
		fmt.Print("Threads参数最大值为2048")
		os.Exit(0)
	}
}

//输出LOGO
func (o *args) printBanner() {
	if len(os.Args) == 1 {
		fmt.Print(o.LOGO)
		fmt.Print(o.USAGE)
		os.Exit(0)
	}
	if o.Help {
		fmt.Print(o.LOGO)
		fmt.Print(o.USAGE)
		fmt.Print(o.HELP)
		os.Exit(0)
	}
	if o.FofaSyntax {
		fmt.Print(o.LOGO)
		fmt.Print(o.USAGE)
		fmt.Print(o.SYNTAX)
		os.Exit(0)
	}
	//打印logo
	fmt.Print(o.LOGO)
}

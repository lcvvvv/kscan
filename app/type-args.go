package app

import (
	"fmt"
	"kscan/lib/sflag"
	"os"
)

type args struct {
	USAGE, HELP, LOGO, SYNTAX string

	Help, Debug, ClosePing, Check, CloseColor, Scan bool
	ScanVersion, DownloadQQwry, CloseCDN            bool
	Output, Proxy, Encoding                         string
	Port, ExcludedPort                              []int
	Path, Host, Target                              []string
	OutputJson, OutputCSV                           string
	Spy, Touch                                      string
	Top, Threads, Timeout                           int
	//hydra模块
	Hydra, HydraUpdate             bool
	HydraUser, HydraPass, HydraMod []string
	//fofa模块
	Fofa                      []string
	FofaField, FofaFixKeyword string
	FofaSize                  int
	FofaSyntax                bool
	//输出修饰
	Match, NotMatch string
}

var Args = args{}

// Parse 初始化参数
func (o *args) Parse() {
	//自定义Usage
	sflag.SetUsage(o.LOGO)
	//定义参数
	o.define()
	//实例化参数值
	sflag.Parse()
	//输出LOGO
	o.printBanner()
}

//定义参数
func (o *args) define() {
	sflag.BoolVar(&o.Help, "h", false)
	sflag.BoolVar(&o.Help, "help", false)
	sflag.BoolVar(&o.Debug, "debug", false)
	sflag.BoolVar(&o.Debug, "d", false)
	//spy模块
	sflag.AutoVarString(&o.Spy, "spy", "None")
	//hydra模块
	sflag.BoolVar(&o.Hydra, "hydra", false)
	sflag.BoolVar(&o.HydraUpdate, "hydra-update", false)
	sflag.StringSpliceVar(&o.HydraUser, "hydra-user")
	sflag.StringSpliceVar(&o.HydraPass, "hydra-pass")
	sflag.StringSpliceVar(&o.HydraMod, "hydra-mod")
	//fofa模块
	sflag.StringSpliceVar(&o.Fofa, "fofa")
	sflag.StringSpliceVar(&o.Fofa, "f")
	sflag.StringVar(&o.FofaField, "fofa-field", "")
	sflag.StringVar(&o.FofaFixKeyword, "fofa-fix-keyword", "")
	sflag.IntVar(&o.FofaSize, "fofa-size", 100)
	sflag.BoolVar(&o.FofaSyntax, "fofa-syntax", false)
	sflag.BoolVar(&o.Scan, "scan", false)
	//kscan模块
	sflag.StringSpliceVar(&o.Target, "target")
	sflag.StringSpliceVar(&o.Target, "t")
	sflag.IntSpliceVar(&o.Port, "port")
	sflag.IntSpliceVar(&o.Port, "p")
	sflag.IntSpliceVar(&o.ExcludedPort, "eP")
	sflag.IntSpliceVar(&o.ExcludedPort, "excluded-port")
	sflag.StringSpliceVar(&o.Path, "path")
	sflag.StringSpliceVar(&o.Host, "host")
	sflag.StringVar(&o.Proxy, "proxy", "")
	sflag.IntVar(&o.Top, "top", 400)
	sflag.IntVar(&o.Threads, "threads", 800)
	sflag.IntVar(&o.Timeout, "timeout", 3)
	sflag.BoolVar(&o.ClosePing, "Pn", false)
	sflag.BoolVar(&o.Check, "check", false)
	sflag.BoolVar(&o.ScanVersion, "sV", false)
	//CDN检测
	sflag.BoolVar(&o.CloseCDN, "Dn", false)
	sflag.BoolVar(&o.DownloadQQwry, "download-qqwry", false)

	//输出模块
	sflag.StringVar(&o.Encoding, "encoding", "utf-8")
	sflag.StringVar(&o.Match, "match", "")
	sflag.StringVar(&o.NotMatch, "not-match", "")
	sflag.StringVar(&o.Output, "o", "")
	sflag.StringVar(&o.Output, "output", "")
	sflag.StringVar(&o.OutputJson, "oJ", "")
	sflag.StringVar(&o.OutputCSV, "oC", "")
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

// CheckArgs 校验参数真实性
func (o *args) CheckArgs() {
	//判断必须的参数是否存在
	if len(o.Target) == 0 && len(o.Fofa) == 0 && o.Spy == "None" && o.DownloadQQwry == false {
		fmt.Print("至少有--target、--fofa、--spy参数中的一个")
		os.Exit(0)
	}
	//判断冲突参数
	if len(o.Target) > 0 && len(o.Fofa) == 0 && o.Spy != "None" && o.Touch == "None" {
		fmt.Print("--target、--fofa、--spy不能同时使用")
		os.Exit(0)
	}
	if len(o.Port) > 0 && o.Top != 400 {
		fmt.Print("--port、--top参数不能同时使用")
		os.Exit(0)
	}
	//判断内容
	if o.Top != 0 && (o.Top > 1000 || o.Top < 1) {
		fmt.Print("TOP参数输入错误,TOP参数应为1-1000之间的整数。")
		os.Exit(0)
	}
	if o.Proxy != "" && sflag.ProxyStrVerification(o.Proxy) {
		fmt.Print("--proxy参数输入错误，其格式应为：http://ip:port，支持socks5/4")
	}
	if o.Threads != 0 && o.Threads > 2048 {
		fmt.Print("--threads参数最大值为2048")
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

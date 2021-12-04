package params

import (
	"fmt"
	"os"

	"kscan/lib/sflag"
)

type OsArgs struct {
	help, debug, scanPing, check, noColor             bool
	target, port, output, proxy, path, host, encoding string
	outputJson                                        string
	USAGE, HELP, LOGO, SYNTAX                         string
	spy, touch                                        string
	top, threads, timeout, rarity                     int
	//hydra模块
	hydra, hydraUpdate             bool
	hydraUser, hydraPass, hydraMod string
	//fofa模块
	fofa, fofaField, fofaFixKeyword string
	fofaSize                        int
	fofaSyntax                      bool
	scan                            bool
}

func (o OsArgs) NoColor() bool {
	return o.noColor
}

func (o OsArgs) Fofa() string {
	return o.fofa
}

func (o OsArgs) Scan() bool {
	return o.scan
}

func (o OsArgs) FofaField() string {
	return o.fofaField
}

func (o OsArgs) FofaFixKeyword() string {
	return o.fofaFixKeyword
}

func (o OsArgs) FofaSize() int {
	return o.fofaSize
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

func (o OsArgs) OutputJson() string {
	return o.outputJson
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

func (o OsArgs) Rarity() int {
	return o.rarity
}

func (o OsArgs) ScanPing() bool {
	return o.scanPing
}

func (o OsArgs) Check() bool {
	return o.check
}

func (o OsArgs) Debug() bool {
	return o.debug
}

func (o OsArgs) Hydra() bool {
	return o.hydra
}

func (o OsArgs) HydraUpdate() bool {
	return o.hydraUpdate
}

func (o OsArgs) HydraMod() string {
	return o.hydraMod
}

func (o OsArgs) HydraUser() string {
	return o.hydraUser
}

func (o OsArgs) HydraPass() string {
	return o.hydraPass
}

func (o OsArgs) Encoding() string {
	return o.encoding
}

func (o OsArgs) Spy() string {
	return o.spy
}

func (o OsArgs) Touch() string {
	return o.touch
}

//初始化参数
func (o *OsArgs) Parse() {
	//自定义Usage
	sflag.SetUsage(o.LOGO)
	//定义参数
	sflag.BoolVar(&o.help, "h", false)
	sflag.BoolVar(&o.help, "help", false)
	sflag.BoolVar(&o.debug, "debug", false)
	sflag.BoolVar(&o.debug, "d", false)
	//spy模块
	sflag.AutoVarString(&o.spy, "spy", "None")
	//touch模块
	sflag.StringVar(&o.touch, "touch", "None")
	//hydra模块
	sflag.BoolVar(&o.hydra, "hydra", false)
	sflag.BoolVar(&o.hydraUpdate, "hydra-update", false)
	sflag.StringVar(&o.hydraUser, "hydra-user", "")
	sflag.StringVar(&o.hydraPass, "hydra-pass", "")
	sflag.StringVar(&o.hydraMod, "hydra-mod", "")
	//fofa模块
	sflag.StringVar(&o.fofa, "fofa", "")
	sflag.StringVar(&o.fofa, "f", "")
	sflag.StringVar(&o.fofaField, "fofa-field", "")
	sflag.StringVar(&o.fofaFixKeyword, "fofa-fix-keyword", "")
	sflag.IntVar(&o.fofaSize, "fofa-size", 100)
	sflag.BoolVar(&o.fofaSyntax, "fofa-syntax", false)
	sflag.BoolVar(&o.scan, "scan", false)
	//kscan模块
	sflag.StringVar(&o.target, "target", "")
	sflag.StringVar(&o.target, "t", "")
	sflag.StringVar(&o.port, "p", "")
	sflag.StringVar(&o.port, "port", "")
	sflag.StringVar(&o.proxy, "proxy", "")
	sflag.StringVar(&o.path, "path", "")
	sflag.StringVar(&o.host, "host", "")
	sflag.IntVar(&o.rarity, "rarity", 9)
	sflag.IntVar(&o.top, "top", 400)
	sflag.IntVar(&o.threads, "threads", 400)
	sflag.IntVar(&o.timeout, "timeout", 3)
	sflag.BoolVar(&o.scanPing, "Pn", false)
	sflag.BoolVar(&o.check, "check", false)
	//输出模块
	sflag.StringVar(&o.encoding, "encoding", "utf-8")
	sflag.StringVar(&o.output, "o", "")
	sflag.StringVar(&o.output, "output", "")
	sflag.StringVar(&o.outputJson, "oJ", "")
	sflag.BoolVar(&o.noColor, "Cn", false)
	//实例化参数值
	sflag.Parse()
}

//初始化函数
func (o *OsArgs) PrintBanner() {
	//不带参数则对应usage
	if len(os.Args) == 1 {
		fmt.Print(o.LOGO)
		fmt.Print(o.USAGE)
		os.Exit(0)
	}
	if o.help {
		fmt.Print(o.LOGO)
		fmt.Print(o.USAGE)
		fmt.Print(o.HELP)
		os.Exit(0)
	}
	if o.fofaSyntax {
		fmt.Print(o.LOGO)
		fmt.Print(o.USAGE)
		fmt.Print(o.SYNTAX)
		os.Exit(0)
	}
	//打印logo
	fmt.Print(o.LOGO)
}

//校验参数真实性
func (o *OsArgs) CheckArgs() {
	//判断必须的参数是否存在
	if o.target == "" && o.fofa == "" && o.spy == "None" && o.touch == "None" {
		fmt.Print("至少有target、fofa、spy、touch参数中的一个")
		os.Exit(0)
	}
	//判断冲突参数
	if o.target != "" && o.fofa != "" && o.spy != "None" && o.touch == "None" {
		fmt.Print("target、fofa、spy、touch不能同时使用")
		os.Exit(0)
	}
	if o.port != "" && o.top != 400 {
		fmt.Print("port、top参数不能同时使用")
		os.Exit(0)
	}
	//判断内容
	if o.port != "" && sflag.MultipleIntVerification(o.port) == false {
		fmt.Print("PORT参数输入错误,其格式应为80，8080，8081-8090")
		os.Exit(0)
	}
	if o.top != 0 && (o.top > 1000 || o.top < 1) {
		fmt.Print("TOP参数输入错误,TOP参数应为1-1000之间的整数。")
		os.Exit(0)
	}
	if o.proxy != "" && sflag.ProxyStrVerification(o.proxy) {
		fmt.Print("PROXY参数输入错误，其格式应为：http://ip:port，支持socks5/4")
	}
	if o.path != "" && sflag.MultipleStrVerification(o.path) {
		fmt.Print("PATH参数输入错误，其格式应为：/asdfasdf，可使用逗号输入多个路径")
	}
	if o.threads != 0 && o.threads > 2048 {
		fmt.Print("Threads参数最大值为2048")
		os.Exit(0)
	}
}

func New(logo string, usage string, help string, syntax string) *OsArgs {
	return &OsArgs{
		LOGO:   logo,
		USAGE:  usage,
		HELP:   help,
		SYNTAX: syntax,
	}
}

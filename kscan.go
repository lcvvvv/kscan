package main

import (
	"embed"
	"fmt"
	"kscan/app"
	"kscan/core/cdn"
	"kscan/core/fofa"
	"kscan/core/hydra"
	"kscan/core/scanner"
	"kscan/core/slog"
	"kscan/core/spy"
	"kscan/core/tips"
	"kscan/lib/color"
	"kscan/lib/misc"
	"kscan/run"
	"os"
	"runtime"
	"time"

	"github.com/lcvvvv/appfinger"
	"github.com/lcvvvv/gonmap"
	"github.com/lcvvvv/pool"
	"github.com/lcvvvv/stdio"
)

//logo信息
const logo = `
     _   __
    /#| /#/ 
    |#|/#/  _____  _____     *     _   _
    |#.#/  /Edge/ /Forum\   /#\   /#\ /#\
    |##|  |#|____ |#|      /Kv2\  |##\|#|
    |#.#\  \r0cky\|#|     /#/_\#\ |#.#.#|
    |#|\#\/\___|#||#|____/#/Rui\#\|#|\##|
    \#| \#\lcvvvv/ \aels/#/ v1.87#\#/ \#/

`

//帮助信息
const help = `
optional arguments:
  -h , --help     show this help message and exit
  -f , --fofa     从fofa获取检测对象，需提前配置环境变量:FOFA_EMAIL、FOFA_KEY
  -t , --target   指定探测对象：
                  IP地址：114.114.114.114
                  IP地址段：114.114.114.114/24,不建议子网掩码小于12
                  IP地址段：114.114.114.114-115.115.115.115
                  URL地址：https://www.baidu.com
                  文件地址：file:/tmp/target.txt
                  剪切板: paste or clipboard
  --spy           网段探测模式，此模式下将自动探测主机可达的内网网段可接收参数为：
                  (空)、192、10、172、all、指定IP地址(将探测该IP地址B段存活网关)
options:
  --check         针对目标地址做指纹识别，仅不会进行端口探测
  --scan          将针对--fofa、--spy提供的目标对象，进行端口扫描和指纹识别
  -p , --port     扫描指定端口，默认会扫描TOP400，支持：80,8080,8088-8090
  -eP, --excluded-port 跳过扫描指定的端口，支持：80,8080,8088-8090
  -o , --output   将扫描结果保存到文件
  -oJ             将扫描结果使用json格式保存到文件
  -oC             将扫描结果使用csv格式保存到文件
  -Pn          	  使用此参数后，将不会进行智能存活性探测，现在默认会开启智能存活性探测，提高效率
  -Cn             使用此参数后，控制台输出结果将不会带颜色。
  -Dn             使用此参数后，将关闭CDN识别功能
  -sV             使用此参数后，将对所有端口进行全探针探测，此参数极度影响效率，慎用！
  --top           扫描经过筛选处理的常见端口TopX，最高支持1000个，默认为TOP400
  --proxy         设置代理(socks5|socks4|https|http)://IP:Port
  --threads       线程参数,默认线程100,最大值为2048
  --path          指定请求访问的目录，只支持单个目录
  --host          指定所有请求的头部Host值
  --timeout       设置超时时间
  --encoding      设置终端输出编码，可指定为：gb2312、utf-8
  --match         对资产返回banner进行检索，剔除不存在关键字的结果记录
  --not-match     对资产返回banner进行检索，剔除存在关键字的结果记录
  --hydra         自动化爆破支持协议：ssh,rdp,ftp,smb,mysql,mssql,oracle,postgresql,mongodb,redis,默认会开启全部
hydra options:
   --hydra-user   自定义hydra爆破用户名:username or user1,user2 or file:username.txt
   --hydra-pass   自定义hydra爆破密码:password or pass1,pass2 or file:password.txt
                  若密码中存在使用逗号的情况，则使用\,进行转义，其他符号无需转义
   --hydra-update 自定义用户名、密码模式，若携带此参数，则为新增模式，会将用户名和密码补充在默认字典后面。否则将替换默认字典。
   --hydra-mod    指定自动化暴力破解模块:rdp or rdp,ssh,smb
fofa options:
   --fofa-syntax  将获取fofa搜索语法说明
   --fofa-size    将设置fofa返回条目数，默认100条
   --fofa-fix-keyword 修饰keyword，该参数中的{}最终会替换成-f参数的值
`

const usage = "usage: kscan [-h,--help,--fofa-syntax] (-t,--target,-f,--fofa,--spy]) [options] [hydra options] [fofa options]\n\n"

const syntax = `title="beijing"			从标题中搜索"北京"			-
header="elastic"		从http头中搜索"elastic"			-
body="网络空间测绘"		从html正文中搜索"网络空间测绘"		-
domain="qq.com"			搜索根域名带有qq.com的网站。		-
icp="京ICP证030173号"		查找备案号为"京ICP证030173号"的网站	搜索网站类型资产
js_name="js/jquery.js"		查找包含js/jquery.js的资产		搜索网站类型资产
js_md5="82ac3f14327a8b7ba49baa208d4eaa15"	查找js源码与之匹配的资产	-
icon_hash="-247388890"		搜索使用此icon的资产。			仅限FOFA高级会员使用
host=".gov.cn"			从url中搜索".gov.cn"			搜索要用host作为名称
port="6379"			查找对应"6379"端口的资产		-
ip="1.1.1.1"			从ip中搜索包含"1.1.1.1"的网站		搜索要用ip作为名称
ip="220.181.111.1/24"		查询IP为"220.181.111.1"的C网段资产	-
status_code="402"		查询服务器状态为"402"的资产		-
protocol="quic"			查询quic协议资产			搜索指定协议类型(在开启端口扫描的情况下有效)
country="CN"			搜索指定国家(编码)的资产。		-
region="Xinjiang"		搜索指定行政区的资产。			-
city="Changsha"			搜索指定城市的资产。			-
cert="baidu"			搜索证书中带有baidu的资产。		-
cert.subject="Oracle"		搜索证书持有者是Oracle的资产		-
cert.issuer="DigiCert"		搜索证书颁发者为DigiCert Inc的资产	-
cert.is_valid=true		验证证书是否有效			仅限FOFA高级会员使用
type=service			搜索所有协议资产			搜索所有协议资产
os="centos"			搜索CentOS资产。			-
server=="Microsoft-IIS"		搜索IIS 10服务器。			-
app="Oracle"			搜索Microsoft-Exchange设备		-
after="2017" && before="2017-10-01"	时间范围段搜索			-
asn="19551"			搜索指定asn的资产。			-
org="Amazon.com, Inc."	搜索指定org(组织)的资产。			-
base_protocol="udp"		搜索指定udp协议的资产。			-
is_fraud=falsenew		排除仿冒/欺诈数据			-
is_honeypot=false		排除蜜罐数据				仅限FOFA高级会员使用
is_ipv6=true			搜索ipv6的资产				搜索ipv6的资产,只接受true和false。
is_domain=true			搜索域名的资产				搜索域名的资产,只接受true和false。
port_size="6"			查询开放端口数量等于"6"的资产		仅限FOFA会员使用
port_size_gt="6"		查询开放端口数量大于"6"的资产		仅限FOFA会员使用
port_size_lt="12"		查询开放端口数量小于"12"的资产		仅限FOFA会员使用
ip_ports="80,161"		搜索同时开放80和161端口的ip		搜索同时开放80和161端口的ip资产(以ip为单位的资产数据)
ip_country="CN"			搜索中国的ip资产。			搜索中国的ip资产
ip_region="Zhejiang"		搜索指定行政区的ip资产。		索指定行政区的资产
ip_city="Hangzhou"		搜索指定城市的ip资产。			搜索指定城市的资产
ip_after="2021-03-18"		搜索2021-03-18以后的ip资产。		搜索2021-03-18以后的ip资产
ip_before="2019-09-09"		搜索2019-09-09以前的ip资产。		搜索2019-09-09以前的ip资产
`

func main() {
	startTime := time.Now()

	//环境初始化
	Init()

	//下载qqwry
	if app.Setting.DownloadQQwry == true {
		slog.Println(slog.INFO, "现在开始下载最新qqwry，请耐心等待！")
		err := cdn.DownloadQQWry()
		if err != nil {
			slog.Println(slog.WARN, "纯真IP库下载失败，请手动下载解压后保存到kscan同一目录")
			slog.Println(slog.WARN, "下载链接： https://qqwry.mirror.noc.one/qqwry.rar")
			slog.Println(slog.WARN, err)
		}
		slog.Println(slog.INFO, "qqwry.dat下载成功！")
		os.Exit(0)
	}

	//spy模块启动
	if app.Setting.Spy != "None" {
		spy.Keyword = app.Setting.Spy
		spy.Scan = app.Setting.Scan
		spy.Start()
		if spy.Scan {
			app.Setting.Target = spy.Target
		}
	}

	//fofa模块初始化
	if len(app.Setting.Fofa) > 0 {
		InitFofa()
		fofa.Run()
		if app.Setting.Check == false && app.Setting.Scan == false {
			slog.Println(slog.WARN, "可以使用--check参数对fofa扫描结果进行存活性及指纹探测，也可以使用--scan参数对fofa扫描结果进行端口扫描")
		}
		if app.Setting.Check == true {
			app.Setting.Target = fofa.GetUrlTarget()
			slog.Println(slog.WARN, "check参数已启用，现在将对fofa扫描结果进行存活性及指纹探测")
		}
		if app.Setting.Scan == true {
			app.Setting.Target = fofa.GetHostTarget()
			slog.Println(slog.WARN, "scan参数已启用，现在将对fofa扫描结果进行端口扫描及指纹探测")
		}
	}
	//Hydra模块初始化
	if app.Setting.Hydra == true {
		slog.Println(slog.INFO, "hydra模块已开启，开始监听暴力破解任务")
		slog.Println(slog.WARN, "当前已开启的hydra模块为：", misc.Intersection(hydra.ProtocolList, app.Setting.HydraMod))
		//加载Hydra模块自定义字典
		hydra.InitCustomAuthMap(app.Setting.HydraUser, app.Setting.HydraPass)
	}
	//kscan模块启动
	if len(app.Setting.Target) > 0 {
		//扫描模块初始化
		InitKscan()
		//开始扫描
		run.Start()
	}
	//计算程序运行时间
	elapsed := time.Since(startTime)
	slog.Printf(slog.INFO, "程序执行总时长为：[%s]", elapsed.String())
	slog.Println(slog.INFO, "若有问题欢迎来我的Github提交Bug[https://github.com/lcvvvv/kscan/]")
}

func Init() {
	app.Args.SetLogo(logo)
	app.Args.SetUsage(usage)
	app.Args.SetHelp(help)
	app.Args.SetSyntax(syntax)
	//参数初始化
	app.Args.Parse()
	//基础输出初始化
	stdio.SetEncoding(app.Args.Encoding)
	//参数合法性校验
	app.Args.CheckArgs()
	//日志初始化
	if app.Args.Debug {
		slog.SetLevel(slog.DEBUG)
	} else {
		slog.SetLevel(slog.INFO)
	}
	//color包初始化
	if os.Getenv("KSCAN_COLOR") == "TRUE" {
		color.Enabled()
	}
	if app.Args.CloseColor == true {
		color.Disabled()
	}
	//pool包初始化
	pool.SetLogger(slog.Debug())
	//配置文件初始化
	app.ConfigInit()
	//Output初始化
	if app.Setting.Output != nil {
		slog.SetOutput(app.Setting.Output)
	}
	fmt.Println("Tips:", tips.GetTips())
	slog.Println(slog.INFO, "当前环境为：", runtime.GOOS, ", 输出编码为：", app.Setting.Encoding)
	if runtime.GOOS == "windows" && app.Setting.CloseColor == true {
		slog.Println(slog.INFO, "在Windows系统下，默认不会开启颜色展示，可以通过添加环境变量开启哦：KSCAN_COLOR=TRUE")
	}
}

//go:embed static/fingerprint.txt
var fingerprintEmbed embed.FS

const (
	qqwryPath       = "qqwry.dat"
	fingerprintPath = "static/fingerprint.txt"
)

func InitKscan() {
	//HTTP指纹库初始化
	fs, _ := fingerprintEmbed.Open(fingerprintPath)
	if n, err := appfinger.InitDatabaseFS(fs); err != nil {
		slog.Println(slog.ERROR, "指纹库加载失败，请检查【fingerprint.txt】文件", err)
	} else {
		slog.Printf(slog.INFO, "成功加载HTTP指纹:[%d]条", n)
	}
	//超时及日志配置
	gonmap.SetLogger(slog.Debug())
	slog.Printf(slog.INFO, "成功加载NMAP探针:[%d]个,指纹[%d]条", gonmap.UsedProbesCount, gonmap.UsedMatchCount)
	//CDN检测初始化
	if app.Setting.CloseCDN == false {
		if _, err := os.Lstat(qqwryPath); os.IsNotExist(err) == true {
			slog.Printf(slog.WARN, "未检测到qqwry.dat,将关闭CDN检测功能，如需开启，请执行kscan --download-qqwry下载该文件")
			app.Setting.CloseCDN = true
		} else {
			slog.Printf(slog.INFO, "检测到qqwry.dat,将自动启动CDN检测功能，可使用-Dn参数关闭该功能")
			scanner.CDNCheck = true
			cdn.Init(qqwryPath)
		}
	}
}

func InitFofa() {
	email := os.Getenv("FOFA_EMAIL")
	key := os.Getenv("FOFA_KEY")
	if email == "" || key == "" {
		slog.Println(slog.WARN, "使用-f/-fofa参数前请先配置环境变量：FOFA_EMAIL、FOFA_KEY")
		slog.Println(slog.ERROR, "如果你是想从文件导入端口扫描任务，请使用-t file:/path/to/file")
	}
	fofa.Init(email, key)
}

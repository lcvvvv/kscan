package scan

import (
	"errors"
	"fmt"
	"kscan/src/lib/misc"
	"kscan/src/lib/net/port"
	"kscan/src/lib/slog"
	"os"
	"regexp"
	"strings"
	"time"
)

type Probes struct {
	//cursor     *Cursor
	commands   []Command
	commandMap map[string]*regexp.Regexp
	probeGroup map[string]*probe
	probeSort  []string
	nowProbe   *probe
	exclude    *port.ProtocolPorts
	target     target
	response   *response
}

//type Cursor struct {
//	Now        int
//	Jump       int
//	IsFallback bool
//	Tab        map[string]int
//}

type Command struct {
	name string
	args string
	//serArgs map[string]interface{}
}

type target struct {
	netloc string
	host   string
	port   int
}

func (this *target) Load(netloc string, host string, port int) {
	this.netloc = netloc
	this.port = port
	this.host = host
}

func (this *Probes) Scan(target target) *response {
	if this.exclude.IsExist(target.port) {
		return nil
	}
	for _, p := range this.probeGroup {
		slog.Info("开始测试:" + p.request.name)
		if !p.Scan(target) {
			//没有返回包，尝试下一种请求包
			continue
		} else {

			this.response = p.response
			fmt.Println(this.response.string)

			if p.Match() {
				fmt.Println(p.response.finger)
				return this.response
			}
			if p.fallback != "" {
				return this.Fallback(p.fallback)
			}
			break
		}
	}
	if this.response.string == "" {
		return nil
	} else {
		return this.response
	}
	//this.cursor.Now = 0
}

func New() *Probes {
	r := Probes{}
	r.response = newResponse()
	r.probeGroup = make(map[string]*probe)
	r.commandMap = map[string]*regexp.Regexp{
		"Exclude": misc.MakeRegexpCompile("^(?:([TU]):)?(\\d+(?:-\\d+)?)$"),
		//Exclude [<protocol>:]<startport>[-<endport>]
		"Probe": misc.MakeRegexpCompile("^(UDP|TCP) ([a-zA-Z0-9-_./]+) (?:q\\|([^|]*)\\|)$"),
		//Probe <protocol> <probename> <probestring>
		"match":  misc.MakeRegexpCompile("^([a-zA-Z0-9-_./]+) m\\|([^|]+)\\|([is]{0,2})(?: (.*))?$"),
		"match=": misc.MakeRegexpCompile("^([a-zA-Z0-9-_./]+) m=([^=]+)=([is]{0,2})(?: (.*))?$"),
		"match%": misc.MakeRegexpCompile("^([a-zA-Z0-9-_./]+) m%([^%]+)%([is]{0,2})(?: (.*))?$"),
		"match@": misc.MakeRegexpCompile("^([a-zA-Z0-9-_./]+) m@([^@]+)@([is]{0,2})(?: (.*))?$"),
		//match <service> <pattern> <patternopt> [<versioninfo>]
		"matchVersioninfoProductname": misc.MakeRegexpCompile("p/([^/]+)/"),
		"matchVersioninfoVersion":     misc.MakeRegexpCompile("v/([^/]+)/"),
		"matchVersioninfoInfo":        misc.MakeRegexpCompile("i/([^/]+)/"),
		"matchVersioninfoHostname":    misc.MakeRegexpCompile("h/([^/]+)/"),
		"matchVersioninfoOS":          misc.MakeRegexpCompile("o/([^/]+)/"),
		"matchVersioninfoDevice":      misc.MakeRegexpCompile("d/([^/]+)/"),
		//  p/vendorproductname/
		//	v/version/
		//	i/info/
		//	h/hostname/
		//	o/operatingsystem/
		//	d/devicetype/

		//  CPE暂时不解析
		//	cpe:/cpename/[a]
		//  cpe:/<part>:<vendor>:<product>:<version>:<update>:<edition>:<language>

		//a for applications,
		//h for hardware platforms, or
		//o for operating systems.

		//$P() 过滤掉不可打印字符
		//$SUBST(1,"_",".")  替换$1中的所有_为.
		//$I(1,">") 16进制转10进制

		//"softmatch": misc.MakeRegexpCompile("^([a-zA-Z0-9-_./]+) m\\|([^|]+)\\|([is]{0,2})"),
		//softmatch <service> <pattern> <patternopt>
		"int": misc.MakeRegexpCompile("^(\\d+)$"),
		"str": misc.MakeRegexpCompile("^([a-zA-Z0-9-_./]+)$"),
	}
	//r.cursor.Now = 0
	//r.cursor.Tab = make(map[string]int)
	return &r
}

func NewTarget() target {
	return target{
		netloc: "",
		host:   "",
		port:   0,
	}
}

func (this *Probes) Load(file *os.File) {
	misc.ReadLineFile(file, this.loadCommand)
	this.probeGroup[this.nowProbe.request.name] = this.nowProbe
	this.probeSort = append(this.probeSort, this.nowProbe.request.name)
	this.nowProbe = nil
	//this.cursor.Now = 0
}

//func (this *Probes) Run(host string, port int) {
//	this.target.netloc = fmt.Sprintf("%s:%d", host, port)
//	this.target.host = host
//	this.target.port = port
//	for this.cursor.Now <= len(this.commands) {
//		this.runCommand(this.commands[this.cursor.Now])
//	}
//}

//func (this *Probes) runCommand(command Command) {
//	switch command.name {
//	case "Exclude":
//		this.exclude = command.serArgs["exclude"].(exclude)
//	case "Probe":
//		this.commands = append(this.commands, command)
//	case "match":
//		command.serArgs = this.loadMatch(command.args)
//		this.commands = append(this.commands, command)
//	case "softmatch":
//
//	case "ports":
//
//	case "sslports":
//
//	case "totalwaitms":
//
//	case "tcpwrappedms":
//
//	case "rarity":
//
//	case "fallback":
//
//	default:
//		return
//	}
//}

func (this *Probes) isCommand(line string) bool {
	var i int
	i = strings.Index(line, "#")
	if i != -1 {
		line = line[:i]
	}
	i = strings.Index(line, " ")
	if i == -1 {
		return false
	}
	return true
}

func (this *Probes) makeCommand(line string) Command {
	i := strings.Index(line, " ")
	command := Command{
		name: line[:i],
		args: line[i+1:],
	}
	commandArr := []string{
		"Exclude", "Probe", "match", "softmatch", "ports", "sslports", "totalwaitms", "tcpwrappedms", "rarity", "fallback",
	}
	if !misc.IsInStrArr(commandArr, command.name) {
		slog.Error(fmt.Sprintf("命令格式不正确，不存在这个命令：%s", command.name))
	}
	//将命令写入command切片
	this.commands = append(this.commands, command)
	return command
}

func (this *Probes) loadCommand(line string) {
	if !this.isCommand(line) {
		return
	}
	command := this.makeCommand(line)
	//fmt.Println(command)
	switch command.name {
	case "Exclude":
		//载入
		this.loadExclude(command.args)
	case "Probe":
		//创建新Probe对象
		p := newProbe()
		//验证是否为第一个对象，若不是则将前一个对象存入切片
		if this.nowProbe != nil {
			this.probeGroup[this.nowProbe.request.name] = this.nowProbe
			this.probeSort = append(this.probeSort, this.nowProbe.request.name)
			//继承属性
			p.tcpwrappedms = this.nowProbe.tcpwrappedms
			p.totalwaitms = this.nowProbe.totalwaitms
		}
		//赋值
		this.nowProbe = p
		//载入
		this.loadProbe(command.args)
	case "match":
		this.loadMatch(command.args, false)
	case "softmatch":
		this.loadMatch(command.args, true)
	case "ports":
		this.loadPorts(command.args, false)
	case "sslports":
		this.loadPorts(command.args, true)
	case "totalwaitms":
		this.nowProbe.totalwaitms = time.Duration(this.makeInt(command.args)) * time.Millisecond
	case "tcpwrappedms":
		this.nowProbe.tcpwrappedms = time.Duration(this.makeInt(command.args)) * time.Millisecond
	case "rarity":
		this.nowProbe.rarity = this.makeInt(command.args)
	case "fallback":
		this.nowProbe.fallback = this.makeStr(command.args)
	}
	//this.cursor.Now++
}

func (this *Probes) makeInt(exprs string) int {
	if !this.commandMap["int"].MatchString(exprs) {
		slog.Error(fmt.Sprintf("数据格式不正确：%s", exprs))
	}
	args := this.commandMap["int"].FindStringSubmatch(exprs)
	return misc.Str2Int(args[1])
}

func (this *Probes) makeStr(exprs string) string {
	if !this.commandMap["str"].MatchString(exprs) {
		slog.Error(fmt.Sprintf("数据格式不正确：%s", exprs))
	}
	args := this.commandMap["str"].FindStringSubmatch(exprs)
	return args[1]
}

func (this *Probes) loadPorts(str string, ssl bool) {
	var p = port.New()
	for _, s := range strings.Split(str, ",") {
		p.Load(s)
	}
	if ssl {
		this.nowProbe.sslports = p
	} else {
		this.nowProbe.ports = p
	}
}

func (this *Probes) loadExclude(str string) {
	var exclude = port.NewProtocolPorts()
	for _, s := range strings.Split(str, ",") {
		r := this.commandMap["Exclude"].FindStringSubmatch(s)
		if r[2] == "" {
			slog.Error("exclude 语句格式错误")
		}
		protocol := r[1]
		exprs := r[2]
		switch protocol {
		case "":
			exclude.TCP.Load(exprs)
			exclude.UDP.Load(exprs)
		case "T":
			exclude.TCP.Load(exprs)
		case "U":
			exclude.UDP.Load(exprs)
		}
	}
	this.exclude = exclude
}

func (this *Probes) loadProbe(s string) {
	//Probe <protocol> <probename> <probestring>
	if !this.commandMap["Probe"].MatchString(s) {
		slog.Error(errors.New("probe 语句格式不正确").Error())
	}
	args := this.commandMap["Probe"].FindStringSubmatch(s)
	if args[1] == "" || args[2] == "" {
		slog.Error(errors.New("probe 参数格式不正确").Error())
	}
	this.nowProbe.request.protocol = args[1]
	this.nowProbe.request.name = args[2]
	this.nowProbe.request.string = args[3]

}

func (this *Probes) loadMatch(s string, soft bool) {
	match := newMatch()
	//"match": misc.MakeRegexpCompile("^([a-zA-Z0-9-_./]+) m\\|([^|]+)\\|([is]{0,2}) (.*)$"),
	//match <service> <pattern>|<patternopt> [<versioninfo>]
	//	"matchVersioninfoProductname": misc.MakeRegexpCompile("p/([^/]+)/"),
	//	"matchVersioninfoVersion":     misc.MakeRegexpCompile("v/([^/]+)/"),
	//	"matchVersioninfoInfo":        misc.MakeRegexpCompile("i/([^/]+)/"),
	//	"matchVersioninfoHostname":    misc.MakeRegexpCompile("h/([^/]+)/"),
	//	"matchVersioninfoOS":          misc.MakeRegexpCompile("o/([^/]+)/"),
	//	"matchVersioninfoDevice":      misc.MakeRegexpCompile("d/([^/]+)/"),
	var args []string
	if this.commandMap["match"].MatchString(s) {
		args = this.commandMap["match"].FindStringSubmatch(s)
	}
	if this.commandMap["match="].MatchString(s) {
		args = this.commandMap["match="].FindStringSubmatch(s)
	}
	if this.commandMap["match%"].MatchString(s) {
		args = this.commandMap["match%"].FindStringSubmatch(s)
	}
	if this.commandMap["match@"].MatchString(s) {
		args = this.commandMap["match@"].FindStringSubmatch(s)
	}
	if args[1] == "" || args[2] == "" {
		slog.Error(errors.New("match 语句参数不正确").Error())
	}
	match.soft = soft
	match.service = args[1]
	match.pattern = args[2]
	match.versioninfo.service = match.service
	match.versioninfo.version = this.getMatchVersionInfo(s, "matchVersioninfoVersion")
	match.versioninfo.productname = this.getMatchVersionInfo(s, "matchVersioninfoProductname")
	match.versioninfo.operatingsystem = this.getMatchVersionInfo(s, "matchVersioninfoOS")
	match.versioninfo.hostname = this.getMatchVersionInfo(s, "matchVersioninfoHostname")
	match.versioninfo.devicetype = this.getMatchVersionInfo(s, "matchVersioninfoDevice")
	match.versioninfo.info = this.getMatchVersionInfo(s, "matchVersioninfoInfo")
	this.nowProbe.matchs = append(this.nowProbe.matchs, match)
}

func (this *Probes) Fallback(ProbeName string) *response {
	this.probeGroup[ProbeName].response = this.response
	if this.probeGroup[ProbeName].Match() {
		return this.probeGroup[ProbeName].response
	} else {
		return this.response
	}
}

func (this *Probes) getMatchVersionInfo(s string, regID string) string {
	if this.commandMap[regID].MatchString(s) {
		return this.commandMap[regID].FindStringSubmatch(s)[1]
	} else {
		return ""
	}
}

func (this *Probes) Show() {
	fmt.Print(this.probeGroup)
}

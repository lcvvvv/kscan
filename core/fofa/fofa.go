package fofa

import (
	"fmt"
	"kscan/app"
	"kscan/core/slog"
	"kscan/lib/color"
	"kscan/lib/fofa"
	"kscan/lib/misc"
	"regexp"
	"strconv"
	"strings"
)

var this *fofa.Fofa
var keywordSlice []string

func Init(email, key string) {
	//设置日志输出器
	fofa.SetLogger(slog.Debug())
	//初始化fofa模块
	this = fofa.New(email, key)
	this.SetSize(app.Setting.FofaSize)
	//获取所有关键字
	keywordSlice = makeKeywordSlice()
}

func Run() {
	//对每个关键字进行查询
	for _, keyword := range keywordSlice {
		slog.Printf(slog.WARN, "本次搜索关键字为：%v", keyword)
		size, results := this.Search(keyword)
		displayResponse(results)
		slog.Printf(slog.INFO, "本次搜索，返回结果总条数为：%d，此次返回条数为：%d", size, len(results))
	}
}

func makeKeywordSlice() []string {
	var keywordSlice []string
	if app.Setting.FofaFixKeyword == "" {
		keywordSlice = app.Setting.Fofa
	} else {
		for _, keyword := range app.Setting.Fofa {
			keyword = strings.ReplaceAll(app.Setting.FofaFixKeyword, "{}", keyword)
			keywordSlice = append(keywordSlice, keyword)
		}
	}
	return keywordSlice
}

func GetUrlTarget() []string {
	var strSlice []string
	for _, result := range this.Results() {
		Fix(&result)
		strSlice = append(strSlice, result.Host)
	}
	strSlice = misc.RemoveDuplicateElement(strSlice)
	return strSlice
}

func GetHostTarget() []string {
	var strSlice []string
	for _, result := range this.Results() {
		strSlice = append(strSlice, result.Ip)
	}
	strSlice = misc.RemoveDuplicateElement(strSlice)
	return strSlice
}

func displayResponse(results []fofa.Result) {
	for _, row := range results {
		Fix(&row)

		m := row.Map()
		m["Header"] = ""
		m["Cert"] = ""
		m["Title"] = ""
		m["Host"] = ""
		m["As_organization"] = ""
		m["Ip"] = ""
		m["Port"] = ""
		m["Country_name"] = ""
		m = misc.FixMap(m)

		if m["Banner"] != "" {
			m["Banner"] = misc.FixLine(m["Banner"])
			m["Banner"] = misc.StrRandomCut(m["Banner"], 20)
		}
		line := fmt.Sprintf("%-30v %-"+strconv.Itoa(misc.AutoWidth(row.Title, 26))+"v %v",
			row.Host,
			row.Title,
			color.StrMapRandomColor(m, app.Setting.CloseColor, []string{"Server"}, []string{}),
		)
		slog.Println(slog.DATA, line)
	}
}

func Fix(r *fofa.Result) {
	if r.Protocol != "" {
		r.Host = fmt.Sprintf("%s://%s:%s", r.Protocol, r.Ip, r.Port)
	}
	if regexp.MustCompile("http([s]?)://.*").MatchString(r.Host) == false && r.Protocol == "" {
		r.Host = "http://" + r.Host
	}
	if r.Title == "" && r.Protocol != "" {
		r.Title = strings.ToUpper(r.Protocol)
	}

	r.Title = misc.FixLine(r.Title)
}

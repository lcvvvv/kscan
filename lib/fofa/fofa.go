package fofa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kscan/app"
	"kscan/lib/color"
	"kscan/lib/misc"
	"kscan/lib/slog"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Fofa struct {
	email, key                     string
	baseUrl, loginPath, searchPath string
	fieldList                      []string

	keywordArr []string
	size       int

	results []Result
}

type ResponseJson struct {
	Error   bool       `json:"error"`
	Mode    string     `json:"mode"`
	Page    int        `json:"page"`
	Query   string     `json:"query"`
	Results [][]string `json:"results"`
	Size    int        `json:"size"`
}

func New(email, key string) *Fofa {
	f := &Fofa{
		email:      email,
		key:        key,
		baseUrl:    "https://fofa.info",
		searchPath: "/api/v1/search/all",
		loginPath:  "/api/v1/info/my",
		fieldList: []string{
			"host", "title", "ip", "domain", "port", "country", "province",
			"city", "country_name", "header", "server", "protocol", "banner",
			"cert", "isp", "as_organization",
		},
	}
	return f
}

func (f *Fofa) LoadArgs() {
	f.loadKeywordArr()
	f.size = app.Setting.FofaSize
}

func (f *Fofa) SearchAll() {
	for _, keyword := range f.keywordArr {
		slog.Warningf("本次搜索关键字为：%v", keyword)

		f.Search(keyword)
	}
}

func (f *Fofa) Search(keyword string) {
	url := f.baseUrl + f.searchPath
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	q := req.URL.Query()
	q.Add("qbase64", misc.Base64Encode(keyword))
	q.Add("email", f.email)
	q.Add("key", f.key)
	q.Add("page", "1")
	q.Add("fields", strings.Join(f.fieldList, ","))
	q.Add("size", misc.Int2Str(f.size))
	q.Add("full", "false")
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		slog.Error(err)
	}
	var responseJson ResponseJson
	if err = json.Unmarshal(body, &responseJson); err != nil {
		slog.Error(body, err)
	}
	r := f.makeResult(responseJson)
	f.results = append(f.results, r...)
	//输出扫描结果
	for _, row := range r {
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
		line := fmt.Sprintf("%-30v %-"+strconv.Itoa(misc.AutoWidth(row.Title, 26))+"v %v\n",
			row.Host,
			row.Title,
			color.StrMapRandomColor(m, app.Setting.CloseColor, []string{"Server"}, []string{}),
		)
		slog.Data(line)
	}
	slog.Infof("本次搜索，返回结果总条数为：%d，此次返回条数为：%d", responseJson.Size, len(responseJson.Results))

	//table.SetPrintColumns(misc.First2UpperForSlice(f.field))
	//t := table.Table(r)
	//fmt.Println(t)
}

func (f *Fofa) makeResult(responseJson ResponseJson) []Result {
	var results []Result
	var result Result

	for _, row := range responseJson.Results {
		m := reflect.ValueOf(&result).Elem()
		for index, f := range f.fieldList {
			f = misc.First2Upper(f)
			m.FieldByName(f).SetString(row[index])
		}
		result.Fix()
		results = append(results, result)
	}
	return results
}

func (f *Fofa) loadKeywordArr() {
	if app.Setting.FofaFixKeyword == "" {
		f.keywordArr = app.Setting.Fofa
	} else {
		for _, keyword := range app.Setting.Fofa {
			keyword = strings.ReplaceAll(app.Setting.FofaFixKeyword, "{}", keyword)
			f.keywordArr = append(f.keywordArr, keyword)
		}
	}
}

func (f *Fofa) Check() {
	var strArr []string
	for _, result := range f.results {
		strArr = append(strArr, result.Host)
	}
	app.Setting.UrlTarget = strArr
}

func (f *Fofa) Scan() {
	var ipArr []string
	var hostArr []string
	for _, result := range f.results {
		ipArr = append(ipArr, result.Ip)
		hostArr = append(hostArr, result.Host)
	}
	ipArr = misc.RemoveDuplicateElement(ipArr)
	hostArr = misc.RemoveDuplicateElement(hostArr)
	app.Setting.HostTarget = ipArr
	app.Setting.UrlTarget = hostArr
}

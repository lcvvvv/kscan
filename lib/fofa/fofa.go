package fofa

import (
	"encoding/json"
	"io/ioutil"
	"kscan/app"
	"kscan/lib/misc"
	"kscan/lib/slog"
	"net/http"
	"strings"
)

type Fofa struct {
	email, key                     string
	baseUrl, loginPath, searchPath string
	fieldList                      []string

	keywordArr []string
	field      []string
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
		baseUrl:    "https://fofa.so",
		searchPath: "/api/v1/search/all",
		loginPath:  "/api/v1/info/my",
		//fieldList: []string{
		//	"host", "title", "ip", "domain", "port", "country", "province", "city", "country_name", "header", "server",
		//	"protocol", "banner", "cert", "isp", "as_number", "as_organization", "latitude", "longitude",
		//},
		fieldList: []string{
			"host", "title", "ip", "domain", "port", "country", "province", "city", "country_name", "header", "server",
			"protocol", "banner", "cert", "isp", "as_organization",
		},
		field: []string{},
	}
	return f
}

func (f *Fofa) LoadArgs() {
	f.loadKeywordArr()
	f.size = app.Setting.FofaSize
	f.loadField()
}

func (f *Fofa) SearchAll() {
	for _, keyword := range f.keywordArr {
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
	//输出扫描结果
	r.Output()
	f.results = append(f.results, r)
}

func (f *Fofa) makeResult(responseJson ResponseJson) Result {
	var r Result
	r.field = f.field
	r.query = responseJson.Query
	r.size = responseJson.Size
	for _, rowArr := range responseJson.Results {
		rowMap := make(map[string]string)
		for i := range misc.Xrange(len(f.fieldList) - 1) {
			rowMap[f.fieldList[i]] = rowArr[i]
		}
		r.result = append(r.result, rowMap)
	}
	return r
}

func (f *Fofa) loadField() {
	for _, field := range app.Setting.FofaField {
		if misc.IsInStrArr(f.fieldList, field) {
			f.field = append(f.field, field)
		}
	}
	if len(f.field) == 0 {
		f.field = []string{"host", "title", "ip", "port"}
	}
}

func (f *Fofa) loadKeywordArr() {
	if app.Setting.FofaFixKeywored == "" {
		f.keywordArr = app.Setting.Fofa
	} else {
		for _, keyword := range app.Setting.Fofa {
			keyword = strings.ReplaceAll(app.Setting.FofaFixKeywored, "{}", keyword)
			f.keywordArr = append(f.keywordArr, keyword)
		}
	}
}

func (f *Fofa) Check() {
	var strArr []string
	for _, result := range f.results {
		for _, row := range result.result {
			strArr = append(strArr, row["host"])
		}
	}
	app.Setting.UrlTarget = strArr
}

func (f *Fofa) Scan() {
	var strArr []string
	for _, result := range f.results {
		for _, row := range result.result {
			strArr = append(strArr, row["ip"])
		}
	}
	strArr = misc.RemoveDuplicateElement(strArr)
	app.Setting.HostTarget = strArr
}

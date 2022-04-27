package fofa

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"kscan/lib/misc"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var logger = Logger(log.New(io.Discard, "", log.Ldate|log.Ltime))

type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}

type Fofa struct {
	email, key                     string
	baseUrl, loginPath, searchPath string
	fieldList                      []string

	size int

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
			"host",
			//"title",
			//"banner",
			//"header",
			"ip", "domain", "port", "country", "province",
			"city", "country_name",
			"server",
			"protocol",
			"cert", "isp", "as_organization",
		},
	}
	return f
}

func (f *Fofa) SetSize(i int) {
	f.size = i
}

func (f *Fofa) Search(keyword string) (int, []Result) {
	url := f.baseUrl + f.searchPath
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	q := req.URL.Query()
	q.Add("qbase64", misc.Base64Encode(keyword))
	q.Add("email", f.email)
	q.Add("key", f.key)
	q.Add("page", "1")
	q.Add("fields", strings.Join(f.fieldList, ","))
	q.Add("size", strconv.Itoa(f.size))
	q.Add("full", "false")
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println(err)
	}
	var responseJson ResponseJson
	if err = json.Unmarshal(body, &responseJson); err != nil {
		logger.Println(body, err)
	}
	r := f.makeResult(responseJson)
	f.results = append(f.results, r...)
	return responseJson.Size, r
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
		results = append(results, result)
	}
	return results
}

func (f *Fofa) Results() []Result {
	return f.results
}

func SetLogger(log Logger) {
	logger = log
}

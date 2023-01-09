package fofa

import (
	"encoding/json"
	"io/ioutil"
	"kscan/lib/misc"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var logger = Logger(log.New(os.Stdout, "[fofa]", log.Ldate|log.Ltime))

type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}

func SetLogger(log Logger) {
	logger = log
}

type Client struct {
	email, key          string
	baseUrl, searchPath string
	fieldList           []string
	size                int
}

type ResponseJson struct {
	Error   bool       `json:"error"`
	Mode    string     `json:"mode"`
	Page    int        `json:"page"`
	Query   string     `json:"query"`
	Results [][]string `json:"results"`
	Size    int        `json:"size"`
}

const (
	baseURL    = "https://fofa.info"
	searchPath = "/api/v1/search/all"
	//loginPath  = "/api/v1/info/my"
)

func New(email, key string) *Client {
	f := &Client{
		email:      email,
		key:        key,
		baseUrl:    baseURL,
		searchPath: searchPath,
		fieldList: []string{
			"host",
			"title",
			"banner",
			"header",
			"ip", "domain", "port", "country", "province",
			"city", "country_name",
			"server",
			"protocol",
			"cert", "isp", "as_organization",
		},
	}
	return f
}

func (f *Client) SetSize(i int) {
	f.size = i
}

func (f *Client) Search(keyword string) (int, []Result) {
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
		return 0, nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println(err)
		return 0, nil
	}
	var responseJson ResponseJson
	if err = json.Unmarshal(body, &responseJson); err != nil {
		logger.Println(body, err)
		return 0, nil
	}
	r := f.makeResult(responseJson)
	return responseJson.Size, r
}

func (f *Client) makeResult(responseJson ResponseJson) (results []Result) {
	for _, row := range responseJson.Results {
		var result Result
		m := reflect.ValueOf(&result).Elem()
		for index, key := range f.fieldList {
			//首字母大写
			key = strings.ToUpper(key[:1]) + key[1:]
			m.FieldByName(key).SetString(row[index])
		}
		results = append(results, result)
	}
	return results
}

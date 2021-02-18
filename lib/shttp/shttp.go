package shttp

import (
	"../../app/config"
	"../../app/params"
	"../misc"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var headerKeys = map[string]string{
	"Server":          "中间件名称",
	"X-Powered-By":    "中间件或开发语言名称",
	"Content-Length":  "返回包长度",
	"Last-Modified":   "最后一次验证日期",
	"Etag":            "连接标签",
	"Accept-Ranges":   "不知道",
	"Date":            "日期",
	"Content-Type":    "正文类型",
	"Set-Cookie":      "设置cookie值",
	"Connection":      "连接类型",
	"Vary":            "不知道",
	"Keep-Alive":      "长链接保存时间",
	"X-Frame-Options": "框架选项",
}

var newHeaderKeys = initNewHeaderKeys()

func initNewHeaderKeys() *os.File {
	if misc.FileIsExist("newHeaderKeys.txt") {
		f, _ := os.Open("newHeaderKeys.txt")
		return f
	} else {
		f, _ := os.Create("newHeaderKeys.txt")
		return f
	}
}

func Get(Url string) (*http.Response, error) {
	request, _ := http.NewRequest("GET", Url, nil)
	request.Header.Add("User-Agent", getUserAgent())
	request.Header.Add("Cookie", "rememberMe=b69375edcb2b3c5084c02bd9690b6625")
	//修改Host头部参数
	if params.SerParams.Host != "" {
		request.Header.Add("Host", params.SerParams.Host)
	}
	tr := &http.Transport{}
	if params.SerParams.Proxy != "" {
		uri, _ := url.Parse(params.SerParams.Proxy)
		(*tr).Proxy = http.ProxyURL(uri)
	}
	(*tr).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	(*tr).DisableKeepAlives = false
	client := &http.Client{}
	client.Transport = tr
	client.Timeout = time.Second * 2
	resp, err := client.Do(request)
	if err == nil {
		//校验http头部
		for key := range resp.Header {
			if headerKeys[key] == "" {
				headerKeys[key] = "New"
				_, _ = newHeaderKeys.WriteString(fmt.Sprintf("%s: %s\n", key, resp.Header.Get(key)))
				fmt.Print("\r", strings.Repeat(" ", 70))
				fmt.Printf("\r[*]发现生僻Http头部：%s: %s\n", key, resp.Header.Get(key))
			}
		}
		//校验http状态码
		if !misc.IsInIntArr(config.Config.HttpCode, resp.StatusCode) {
			resp = nil
			err = errors.New("HttpStatusCode不在范围内")
			return resp, err
		}
		//检测内容是否压缩
		//if resp.Header.Get("Content-Encoding") == "gzip" {
		//	var reader io.ReadCloser
		//	reader, err = gzip.NewReader(resp.Body)
		//	if err != nil {
		//		resp.Body = reader
		//	}
		//}
	}
	return resp, err
}

func getUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(config.Config.UserAgents))
	return config.Config.UserAgents[i]
}

func Header2String(header http.Header) string {
	var result string
	for i := range header {
		result += fmt.Sprintf("%s: %s\n", i, header.Get(i))
	}
	return result
}

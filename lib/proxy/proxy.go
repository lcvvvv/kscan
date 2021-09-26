package proxy

import (
	"context"
	"crypto/tls"
	"kscan/lib/misc"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Check(Url string) bool {
	request, err := http.NewRequest("GET", "http://www.baidu.com", nil)
	if err != nil {
		return false
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)")
	request.Close = true
	tr := &http.Transport{}
	(*tr).TLSClientConfig = &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS10}
	(*tr).DisableKeepAlives = false
	client := &http.Client{}
	//修改HTTP超时时间
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	request.WithContext(ctx)
	client.Timeout = 5 * time.Second
	//修改代理选项
	uri, _ := url.Parse(Url)
	(*tr).Proxy = http.ProxyURL(uri)
	client.Transport = tr
	//发送数据包
	resp, err := client.Do(request)
	if err != nil {
		return false
	}
	//确认返回包是否正确
	buf := misc.ReadAll(resp.Body, 5*time.Second)
	result := string(buf)
	if strings.Contains(result, "baidu") == false {
		return false
	}
	return true
}

package shttp

import (
	"../../app/config"
	"../../app/params"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func Get(Url string) (*http.Response, error) {
	request, _ := http.NewRequest("GET", Url, nil)
	request.Header.Add("User-Agent", getUserAgent())
	request.Header.Add("Cookie", "rememberMe=b69375edcb2b3c5084c02bd9690b6625")
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
		result += fmt.Sprintf("%d: %d\n", i, header.Get(i))
	}
	return result
}

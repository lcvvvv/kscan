package run

import (
	"fmt"
	"github.com/lcvvvv/urlparse"
	"kscan/app"
	"kscan/lib/gonmap"
	"kscan/lib/shttp"
	"kscan/lib/slog"
	"strings"
	"time"
)

func GetPortBanner(expr string, nmap *gonmap.Nmap) *PortInformation {
	u, _ := urlparse.Load(expr)
	r := NewPortInformation(u)

	//如果未指定协议类型则进行协议检测
	if r.Target.Scheme != "http" && r.Target.Scheme != "https" {
		r.LoadGonmapPortInformation(nmap.SafeScan(u.Netloc, u.Port, time.Second*120))
		if r.ErrorMsg != nil {
			slog.Debug(r.ErrorMsg.Error())
			for i := 0; i < 3; i++ {
				if strings.Contains(r.ErrorMsg.Error(), "STEP3:READ") {
					r.LoadGonmapPortInformation(nmap.SafeScan(u.Netloc, u.Port, time.Second*120))
				}
				if r.ErrorMsg == nil {
					break
				}
			}
		}
		slog.Debug(fmt.Sprintln(u.UnParse(), r.ErrorMsg, r.Target.Scheme, r.Finger.Service))
	}
	//如果端口关闭则直接返回
	if r.Status == "CLOSED" {
		return r
	}

	//如果协议类型为HTTP协议，则进行HTTPbanner识别
	if r.Target.Scheme == "http" {
		r.Target.Scheme = "http"
		r.Target.Path = app.CConfig.Path
		r.LoadHttpFinger(getUrlBanner(r.Target))
	}
	if r.Target.Scheme == "ssl" || r.Target.Scheme == "https" {
		r.Target.Scheme = "https"
		r.Target.Path = app.CConfig.Path
		r.LoadHttpFinger(getUrlBanner(r.Target))
	}
	return r
}

func getUrlBanner(url *urlparse.URL) *HttpFinger {
	r := NewHttpFinger()
	resp, err := shttp.Get(url.UnParse())
	if err != nil {
		if strings.Contains(err.Error(), "server gave HTTP response") {
			//HTTP协议重新获取指纹
			url.Scheme = "http"
			return getUrlBanner(url)
		}
		if strings.Contains(err.Error(), "malformed HTTP response") {
			//HTTP协议重新获取指纹
			url.Scheme = "https"
			return getUrlBanner(url)
		}
		slog.Debug(err.Error())
		return r
	}
	if url.Scheme == "https" {
		r.LoadHttpsResponse(resp)
	}
	r.LoadHttpResponse(url, resp)
	r.MakeInfo()
	//res.Info = makeResultInfo(res)
	return r
}

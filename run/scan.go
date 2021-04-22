package run

import (
	"fmt"
	"github.com/go-ping/ping"
	"github.com/lcvvvv/gonmap"
	"github.com/lcvvvv/urlparse"
	"kscan/app"
	"kscan/lib/shttp"
	"kscan/lib/slog"
	"strings"
	"time"
)

func GetPortBanner(expr string, nmap *gonmap.Nmap) *PortInformation {
	u, _ := urlparse.Load(expr)
	r := NewPortInformation(u)
	if app.Config.PingAliveMap != nil {
		if value, _ := app.Config.PingAliveMap.Load(u.Netloc); value.(int) == 0 {
			app.Config.PingAliveMap.Store(u.Netloc, PingAlive(u.Netloc))
		}
		if value, _ := app.Config.PingAliveMap.Load(u.Netloc); value.(int) == -1 {
			return r.CLOSED()
		}
	}
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
	if r.Status == "CLOSED" {
		return r
	}
	if r.Finger.Service == "http" {
		r.Target.Scheme = "http"
		r.Target.Path = app.Config.Path
		r.LoadHttpFinger(getUrlBanner(r.Target))
	}
	if r.Finger.Service == "ssl" || r.Finger.Service == "https" {
		r.Target.Scheme = "https"
		r.Target.Path = app.Config.Path
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

func PingAlive(ip string) int {
	p, err := ping.NewPinger(ip)
	if err != nil {
		slog.Debug(err.Error())
		return -1
	}
	p.Count = 2
	p.Timeout = time.Second * 3
	err = p.Run() // Blocks until finished.
	if err != nil {
		slog.Debug(err.Error())
	}
	s := p.Statistics()
	if s.PacketsRecv > 0 {
		return 1
	}
	return -1
}

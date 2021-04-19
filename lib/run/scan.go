package run

import (
	"fmt"
	"github.com/go-ping/ping"
	"github.com/lcvvvv/gonmap"
	"github.com/lcvvvv/urlparse"
	"kscan/app"
	"kscan/lib/shttp"
	"kscan/lib/slog"
	"strconv"
	"strings"
	"time"
)

func GetPortBanner(expr string, nmap *gonmap.Nmap) *PortInformation {
	u, _ := urlparse.Load(expr)
	r := NewPortInformation(u)
	if app.Config.PingAliveMap != nil {
		if _, ok := app.Config.PingAliveMap[u.Netloc]; ok != true {
			app.Config.PingAliveMap[u.Netloc] = PingAlive(u.Netloc)
		}
		slog.Debug(u.Netloc + "ping探测结果为：" + strconv.FormatBool(app.Config.PingAliveMap[u.Netloc]))
		if app.Config.PingAliveMap[u.Netloc] != true {
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

//func getTcpBanner(s string) portInfo {
//	var res portInfo
//	url, _ := urlparse.Load(s)
//	res.Url = s
//	res.Netloc = url.Netloc
//	res.Portid = url.Port
//	res.Protocol = getProtocol(s)
//	conn, err := net.DialTimeout("tcp", s, time.Second*app.Config.Timeout)
//	if err != nil {
//		res.Alive = false
//		res.Banner = ""
//		if strings.Contains(err.Error(), "too many") {
//			//发现存在线程过高错误
//			slog.Errorf("当前线程过高，请降低线程!或者请执行\"ulimit -n 50000\"命令放开操作系统限制,MAC系统可能还需要执行：\"launchctl limit maxfiles 50000 50000\"")
//		}
//		slog.Debug(err.Error())
//	} else {
//		_ = conn.SetReadDeadline(time.Now().Add(time.Second * app.Config.Timeout))
//		res.Alive = true
//		res.KeywordFinger.errorMsg = errors.New("非Web端口")
//		res.HashFinger.errorMsg = errors.New("非Web端口")
//		_, _ = conn.Write([]byte("test\r\n"))
//		Bytes := make([]byte, 1024)
//		i, _ := conn.Read(Bytes)
//		res.Banner = string(Bytes[:i])
//		res.Banner = misc.FixLine(res.Banner)
//		conn.Close()
//	}
//	return res
//}

func PingAlive(ip string) bool {
	p, err := ping.NewPinger(ip)
	if err != nil {
		slog.Debug(err.Error())
	}
	p.Count = 2
	p.Timeout = time.Second * 3
	err = p.Run() // Blocks until finished.
	if err != nil {
		slog.Debug(err.Error())
	}
	s := p.Statistics()
	if s.PacketsRecv > 0 {
		return true
	}
	return false
}

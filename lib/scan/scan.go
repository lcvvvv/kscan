package scan

import (
	"fmt"
	"github.com/lcvvvv/gonmap"
	"github.com/lcvvvv/urlparse"
	"kscan/app"
	"kscan/lib/shttp"
	"kscan/lib/slog"
	"strings"
)

func GetPortBanner(expr string, nmap *gonmap.Nmap) *PortInformation {
	u, _ := urlparse.Load(expr)
	r := NewPortInformation(u)
	r.LoadGonmapPortInformation(nmap.Scan(u.Netloc, u.Port))

	if r.Status == "CLOSED" {
		return r
	}
	fmt.Println(r.Finger.Service)
	if r.Finger.Service == "http" {
		r.Target.Scheme = "http"
		r.Target.Path = app.Config.Path
		r.LoadHttpFinger(getUrlBanner(r.Target))
	}
	if r.Finger.Service == "ssl" {
		r.Target.Scheme = "https"
		r.Target.Path = app.Config.Path
		r.LoadHttpFinger(getUrlBanner(r.Target))
	}
	return r
}

func getUrlBanner(url *urlparse.URL) *HttpFinger {
	r := NewHttpFinger()
	fmt.Println(url.UnParse())
	resp, err := shttp.Get(url.UnParse())
	if err != nil {
		if strings.Contains(err.Error(), "too many") {
			//发现存在线程过高错误
			slog.Errorf("当前线程过高，请降低线程!或者请执行\"ulimit -n 50000\"命令放开操作系统限制,MAC系统可能还需要执行：\"launchctl limit maxfiles 50000 50000\"")
		}
		if strings.Contains(err.Error(), "server gave HTTP response") {
			//HTTP协议重新获取指纹
			url.Netloc = "http"
			return getUrlBanner(url)
		}
		slog.Debug(err.Error())
		return r
	}
	r.LoadHttpResponse(url, resp)
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

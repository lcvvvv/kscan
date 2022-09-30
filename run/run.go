package run

import (
	"fmt"
	"github.com/lcvvvv/appfinger"
	"github.com/lcvvvv/gonmap"
	"github.com/lcvvvv/simplehttp"
	"kscan/app"
	"kscan/core/cdn"
	"kscan/core/hydra"
	"kscan/core/scanner"
	"kscan/core/slog"
	"kscan/lib/color"
	"kscan/lib/misc"
	"kscan/lib/uri"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Start() {
	DomainScanner = generateDomainScanner()
	IPScanner = generateIPScanner()
	PortScanner = generatePortScanner()
	URLScanner = generateURLScanner()
	HydraScanner = generateHydraScanner()

	//扫描器进入监听状态
	start()
	//下发扫描任务
	var wg = &sync.WaitGroup{}
	go watchDog(wg)
	for _, expr := range app.Setting.Target {
		pushTarget(expr)
	}
	slog.Println(slog.INFO, "所有扫描任务已下发完毕")
	wg.Wait()
}

func pushTarget(expr string) {
	if expr == "" {
		return
	}
	if uri.IsIP(expr) {
		IPScanner.Push(net.ParseIP(expr))
		if app.Setting.Check == true {
			pushURLTarget(uri.URLParse("http://"+expr), nil)
			pushURLTarget(uri.URLParse("https://"+expr), nil)
		}
		return
	}
	if uri.IsCIDR(expr) {
		for _, ip := range uri.CIDRToIP(expr) {
			pushTarget(ip.String())
		}
		return
	}
	if uri.IsIPRanger(expr) {
		for _, ip := range uri.RangerToIP(expr) {
			pushTarget(ip.String())
		}
		return
	}
	if uri.IsDomain(expr) {
		DomainScanner.Push(expr)
		pushURLTarget(uri.URLParse("http://"+expr), nil)
		pushURLTarget(uri.URLParse("https://"+expr), nil)
		return
	}
	if uri.IsHostPath(expr) {
		pushTarget(uri.GetNetlocWithHostPath(expr))
		pushURLTarget(uri.URLParse("http://"+expr), nil)
		pushURLTarget(uri.URLParse("https://"+expr), nil)
		return
	}
	if uri.IsNetlocPort(expr) {
		netloc, port := uri.SplitWithNetlocPort(expr)
		pushTarget(netloc)
		if uri.IsIP(netloc) {
			PortScanner.Push(net.ParseIP(netloc), port)
		}
		if uri.IsDomain(netloc) {
			pushURLTarget(uri.URLParse("http://"+expr), nil)
			pushURLTarget(uri.URLParse("https://"+expr), nil)
		}
		return
	}
	if uri.IsURL(expr) {
		pushURLTarget(uri.URLParse(expr), nil)
		pushTarget(uri.GetNetlocWithURL(expr))
		return
	}
	slog.Println(slog.WARN, "无法识别的Target字符串:", expr)
}

var (
	DomainScanner *scanner.DomainClient
	IPScanner     *scanner.IPClient
	PortScanner   *scanner.PortClient
	URLScanner    *scanner.URLClient
	HydraScanner  *scanner.HydraClient
)

func start() {
	go DomainScanner.Start()
	go IPScanner.Start()
	go PortScanner.Start()
	go URLScanner.Start()
	go HydraScanner.Start()
}

func stop() {
	DomainScanner.Stop()
	IPScanner.Stop()
	PortScanner.Stop()
	URLScanner.Stop()
	HydraScanner.Stop()
}

func generateDomainScanner() *scanner.DomainClient {
	DomainConfig := scanner.DefaultConfig()
	DomainConfig.Threads = 10
	client := scanner.NewDomainScanner(DomainConfig)
	client.HandlerRealIP = func(domain string, ip net.IP) {
		IPScanner.Push(ip)
	}
	client.HandlerIsCDN = func(domain, CDNInfo string) {
		outputCDNRecord(domain, CDNInfo)
	}
	client.HandlerError = func(domain string, err error) {
		slog.Println(slog.DEBUG, "DomainScanner Error: ", domain, err)
	}
	return client
}

func generateIPScanner() *scanner.IPClient {
	IPConfig := scanner.DefaultConfig()
	IPConfig.Threads = 200
	IPConfig.Timeout = 200 * time.Millisecond
	client := scanner.NewIPScanner(IPConfig)
	client.HandlerDie = func(addr net.IP) {
		slog.Println(slog.DEBUG, addr.String(), " is die")
	}
	client.HandlerAlive = func(addr net.IP) {
		//启用端口存活性探测任务下发器
		for _, port := range app.Setting.Port {
			PortScanner.Push(addr, port)
		}
	}
	client.HandlerError = func(addr net.IP, err error) {
		slog.Println(slog.DEBUG, "IPScanner Error: ", addr.String(), err)
	}
	return client
}

func generatePortScanner() *scanner.PortClient {
	PortConfig := scanner.DefaultConfig()
	PortConfig.Threads = app.Setting.Threads
	PortConfig.Timeout = time.Millisecond * 500
	if app.Setting.ScanVersion == true {
		PortConfig.DeepInspection = true
	}
	client := scanner.NewPortScanner(PortConfig)
	client.HandlerClosed = func(addr net.IP, port int) {
		//nothing
	}
	client.HandlerOpen = func(addr net.IP, port int) {
		//nothing
	}
	client.HandlerNotMatched = func(addr net.IP, port int, response string) {
		outputUnknownResponse(addr, port, response)
	}
	client.HandlerMatched = func(addr net.IP, port int, response *gonmap.Response) {
		URLRaw := fmt.Sprintf("%s://%s:%d", response.FingerPrint.Service, addr.String(), port)
		URL, _ := url.Parse(URLRaw)
		if appfinger.SupportCheck(URL.Scheme) == true {
			pushURLTarget(URL, response)
			return
		}
		if app.Setting.Hydra == true {
			if protocol := response.FingerPrint.Service; hydra.Ok(protocol) {
				HydraScanner.Push(addr, port, protocol)
			}
		}
		outputNmapFinger(URL, response)
	}
	client.HandlerError = func(addr net.IP, port int, err error) {
		slog.Println(slog.DEBUG, "PortScanner Error: ", fmt.Sprintf("%s:%d", addr.String(), port), err)
	}
	return client
}

func generateURLScanner() *scanner.URLClient {
	URLConfig := scanner.DefaultConfig()
	URLConfig.Threads = app.Setting.Threads/2 + 1

	client := scanner.NewURLScanner(URLConfig)
	client.HandlerMatched = func(URL *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint) {
		outputAppFinger(URL, banner, finger)
	}
	client.HandlerError = func(url *url.URL, err error) {
		slog.Println(slog.DEBUG, "URLScanner Error: ", url.String(), err)
	}
	return client
}

func generateHydraScanner() *scanner.HydraClient {
	HydraConfig := scanner.DefaultConfig()
	HydraConfig.Threads = 10

	client := scanner.NewHydraScanner(HydraConfig)
	client.HandlerSuccess = func(addr net.IP, port int, protocol string, auth *hydra.Auth) {
		outputHydraSuccess(addr, port, protocol, auth)
	}
	client.HandlerError = func(addr net.IP, port int, protocol string, err error) {
		slog.Println(slog.DEBUG, fmt.Sprintf("%s://%s:%d", protocol, addr.String(), port), err)
	}
	return client
}

func outputHydraSuccess(addr net.IP, port int, protocol string, auth *hydra.Auth) {
	var target = fmt.Sprintf("%s://%s:%d", protocol, addr.String(), port)
	var m = auth.Map()
	outputHandler(target, color.Important("CrackSuccess"), m)
}

func outputNmapFinger(URL *url.URL, resp *gonmap.Response) {
	if responseFilter(resp.Raw) == true {
		return
	}
	finger := resp.FingerPrint
	m := misc.ToMap(finger)
	m["Response"] = resp.Raw
	//补充归属地信息
	if app.Setting.CloseCDN == false {
		result, _ := cdn.Find(URL.Hostname())
		m["Addr"] = result
	}
	outputHandler(URL.String(), finger.Service, m)
}

func outputAppFinger(URL *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint) {
	if responseFilter(banner.Response, banner.Cert) == true {
		return
	}
	m := misc.ToMap(finger)
	//补充归属地信息
	if app.Setting.CloseCDN == false {
		result, _ := cdn.Find(URL.Hostname())
		m["Addr"] = result
	}
	m["Service"] = URL.Scheme
	m["FoundDomain"] = banner.FoundDomain
	m["FoundIP"] = banner.FoundIP
	m["FingerPrint"] = m["ProductName"]
	m["Response"] = banner.Response
	m["Cert"] = banner.Cert
	m["Header"] = banner.Header
	m["Body"] = banner.Body
	m["ICP"] = banner.ICP
	delete(m, "ProductName")
	outputHandler(URL.String(), banner.Title, m)

}

func outputCDNRecord(domain, info string) {
	if responseFilter(info) == true {
		return
	}
	//输出结果
	domain = fmt.Sprintf("cdn://%s", domain)
	outputHandler(domain, "CDN资产", map[string]string{"CDNInfo": info})
}

func outputUnknownResponse(addr net.IP, port int, response string) {
	if responseFilter(response) == true {
		return
	}
	//输出结果
	target := fmt.Sprintf("unknown://%s:%d", addr.String(), port)
	outputHandler(target, "无法识别该协议", map[string]string{"Response": response})
}

func responseFilter(strArgs ...string) bool {
	if app.Setting.Match == "" {
		return false
	}

	for _, str := range strArgs {
		if strings.Contains(str, app.Setting.Match) == false {
			return true
		}
		//if strings.Contains(str, app.Setting.NotMatch) == true {
		//	return true
		//}
	}
	return false
}

var (
	disableKey       = []string{"MatchRegexString", "Service", "ProbeName", "Response", "Cert", "Header", "Body"}
	importantKey     = []string{"ProductName", "DeviceType"}
	varyImportantKey = []string{"Hostname", "FingerPrint", "ICP"}
)

func getHTTPDigest(s string) string {
	var length = 24
	var digestBuf []rune
	for _, r := range []rune(s) {
		buf := []byte(string(r))
		if len(digestBuf) == length {
			return string(digestBuf)
		}
		if len(buf) > 1 {
			digestBuf = append(digestBuf, r)
		}
	}
	return string(digestBuf) + misc.StrRandomCut(s, length-len(digestBuf))
}

func getRawDigest(s string) string {
	var length = 24
	if len(s) < length {
		return s
	}
	var digestBuf []rune
	for _, r := range []rune(s) {
		if len(digestBuf) == length {
			return string(digestBuf)
		}
		if 0x20 <= r && r <= 0x7E {
			digestBuf = append(digestBuf, r)
		}
	}
	return string(digestBuf) + misc.StrRandomCut(s, length-len(digestBuf))
}

func outputHandler(target, keyword string, m map[string]string) {
	m = misc.FixMap(m)
	if respRaw := m["Response"]; respRaw != "" {
		if m["Service"] == "http" || m["Service"] == "https" {
			m["Digest"] = strconv.Quote(getHTTPDigest(respRaw))
		} else {
			m["Digest"] = strconv.Quote(getRawDigest(respRaw))
		}
	}
	sourceMap := misc.CloneMap(m)
	for _, keyword := range disableKey {
		delete(m, keyword)
	}
	for key, value := range m {
		if key == "FingerPrint" {
			continue
		}
		m[key] = misc.StrRandomCut(value, 24)
	}
	fingerPrint := color.StrMapRandomColor(m, true, importantKey, varyImportantKey)
	fingerPrint = misc.FixLine(fingerPrint)
	format := "%-30v %-" + strconv.Itoa(misc.AutoWidth(color.Clear(keyword), 26+color.Count(keyword))) + "v %s"
	printStr := fmt.Sprintf(format, target, keyword, fingerPrint)
	slog.Println(slog.DATA, printStr)

	if jw := app.Setting.OutputJson; jw != nil {
		m["target"] = target
		m["keyword"] = keyword
		jw.Push(sourceMap)
	}
}

func pushURLTarget(URL *url.URL, response *gonmap.Response) {
	var cli *http.Client
	//判断是否初始化client
	if app.Setting.Proxy != "" || app.Setting.Timeout != 3*time.Second {
		cli = simplehttp.NewClient()
	}
	//判断是否需要设置代理
	if app.Setting.Proxy != "" {
		simplehttp.SetProxy(cli, app.Setting.Proxy)
	}
	//判断是否需要设置超时参数
	if app.Setting.Timeout != 3*time.Second {
		simplehttp.SetTimeout(cli, app.Setting.Timeout)
	}

	//判断是否存在请求修饰性参数
	if len(app.Setting.Host) == 0 && len(app.Setting.Path) == 0 {
		URLScanner.Push(URL, response, nil, cli)
		return
	}

	//如果存在，则逐一建立请求下发队列
	var reqs []*http.Request
	for _, host := range app.Setting.Host {
		URL.Host = host
		req, _ := simplehttp.NewRequest(http.MethodGet, URL.String(), nil)
		reqs = append(reqs, req)
	}
	for _, path := range app.Setting.Path {
		URL.Path = path
		req, _ := simplehttp.NewRequest(http.MethodGet, URL.String(), nil)
		reqs = append(reqs, req)
	}
	for _, req := range reqs {
		URLScanner.Push(req.URL, response, req, cli)
	}
}

func watchDog(wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		time.Sleep(time.Second * 1)
		var (
			nDomain = DomainScanner.RunningThreads()
			nIP     = IPScanner.RunningThreads()
			nPort   = PortScanner.RunningThreads()
			nURL    = URLScanner.RunningThreads()
			nHydra  = HydraScanner.RunningThreads()
		)
		if time.Now().Unix()%180 == 0 {
			warn := fmt.Sprintf("当前存活协程数：Domain：%d 个，IP：%d 个，Port：%d 个，URL：%d 个，Hydra：%d 个", nDomain, nIP, nPort, nURL, nHydra)
			slog.Println(slog.WARN, warn)
		}
		if nDomain != 0 {
			continue
		}

		if nIP != 0 {
			continue
		}
		if nPort != 0 {
			continue
		}
		if nURL != 0 {
			continue
		}
		if nHydra != 0 {
			continue
		}
		break
	}
	stop()
	wg.Done()
}

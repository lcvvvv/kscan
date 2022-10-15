package run

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/lcvvvv/appfinger"
	"github.com/lcvvvv/gonmap"
	"github.com/lcvvvv/simplehttp"
	"github.com/lcvvvv/stdio/chinese"
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
	"runtime"
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
	stop()
}

func pushTarget(expr string) {
	if expr == "" {
		return
	}
	if expr == "paste" || expr == "clipboard" {
		if clipboard.Unsupported == true {
			slog.Println(slog.ERROR, runtime.GOOS, "clipboard unsupported")
		}
		clipboardStr, _ := clipboard.ReadAll()
		for _, line := range strings.Split(clipboardStr, "\n") {
			line = strings.ReplaceAll(line, "\r", "")
			pushTarget(line)
		}
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
		pushURLTarget(uri.URLParse("http://"+expr), nil)
		pushURLTarget(uri.URLParse("https://"+expr), nil)
		if app.Setting.Check == false {
			pushTarget(uri.GetNetlocWithHostPath(expr))
		}
		return
	}
	if uri.IsNetlocPort(expr) {
		netloc, port := uri.SplitWithNetlocPort(expr)
		if uri.IsIP(netloc) {
			PortScanner.Push(net.ParseIP(netloc), port)
		}
		if uri.IsDomain(netloc) {
			pushURLTarget(uri.URLParse("http://"+expr), nil)
			pushURLTarget(uri.URLParse("https://"+expr), nil)
		}
		if app.Setting.Check == false {
			pushTarget(netloc)
		}
		return
	}
	if uri.IsURL(expr) {
		pushURLTarget(uri.URLParse(expr), nil)
		if app.Setting.Check == false {
			pushTarget(uri.GetNetlocWithURL(expr))
		}
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
	IPConfig.HostDiscoverClosed = app.Setting.ClosePing
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

func getTimeout(i int) time.Duration {
	switch {
	case i > 10000:
		return time.Millisecond * 200
	case i > 5000:
		return time.Millisecond * 300
	case i > 1000:
		return time.Millisecond * 400
	default:
		return time.Millisecond * 500
	}
}

func generatePortScanner() *scanner.PortClient {
	PortConfig := scanner.DefaultConfig()
	PortConfig.Threads = app.Setting.Threads
	PortConfig.Timeout = getTimeout(len(app.Setting.Port))
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
	URL, _ := url.Parse(target)
	outputHandler(URL, color.Important("CrackSuccess"), m)
}

func outputNmapFinger(URL *url.URL, resp *gonmap.Response) {
	if responseFilter(resp.Raw) == true {
		return
	}
	finger := resp.FingerPrint
	m := misc.ToMap(finger)
	m["Response"] = resp.Raw
	m["IP"] = URL.Hostname()
	m["Port"] = URL.Port()
	//补充归属地信息
	if app.Setting.CloseCDN == false {
		result, _ := cdn.Find(URL.Hostname())
		m["Addr"] = result
	}
	outputHandler(URL, finger.Service, m)
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
	m["Response"] = banner.Response
	m["Cert"] = banner.Cert
	m["Header"] = banner.Header
	m["Body"] = banner.Body
	m["ICP"] = banner.ICP
	m["FingerPrint"] = m["ProductName"]
	delete(m, "ProductName")
	//增加IP、Domain、Port字段
	m["Port"] = uri.GetURLPort(URL)
	if m["Port"] == "" {
		slog.Println(slog.WARN, "无法获取端口号：", URL)
	}
	if hostname := URL.Hostname(); uri.IsIP(hostname) {
		m["IP"] = hostname
	} else {
		m["Domain"] = hostname
		if v, ok := scanner.DomainDatabase.Load(hostname); ok {
			m["IP"] = v.(string)
		}
	}
	outputHandler(URL, banner.Title, m)
}

func outputCDNRecord(domain, info string) {
	if responseFilter(info) == true {
		return
	}
	//输出结果
	target := fmt.Sprintf("cdn://%s", domain)
	URL, _ := url.Parse(target)
	outputHandler(URL, "CDN资产", map[string]string{
		"CDNInfo": info,
		"Domain":  domain,
	})
}

func outputUnknownResponse(addr net.IP, port int, response string) {
	if responseFilter(response) == true {
		return
	}
	//输出结果
	target := fmt.Sprintf("unknown://%s:%d", addr.String(), port)
	URL, _ := url.Parse(target)
	outputHandler(URL, "无法识别该协议", map[string]string{
		"Response": response,
		"IP":       URL.Hostname(),
		"Port":     strconv.Itoa(port),
	})
}

func responseFilter(strArgs ...string) bool {
	var match = app.Setting.Match
	var notMatch = app.Setting.NotMatch

	if match != "" {
		for _, str := range strArgs {
			//主要结果中包含关键则，则会显示
			if strings.Contains(str, app.Setting.Match) == true {
				return false
			}
		}
	}

	if notMatch != "" {
		for _, str := range strArgs {
			//主要结果中包含关键则，则会显示
			if strings.Contains(str, app.Setting.NotMatch) == true {
				return true
			}
		}
	}
	return false
}

var (
	disableKey       = []string{"MatchRegexString", "Service", "ProbeName", "Response", "Cert", "Header", "Body", "IP"}
	importantKey     = []string{"ProductName", "DeviceType"}
	varyImportantKey = []string{"Hostname", "FingerPrint", "ICP"}
)

func getHTTPDigest(s string) string {
	var length = 24
	var digestBuf []rune
	_, body := simplehttp.SplitHeaderAndBody(s)
	body = chinese.ToUTF8(body)
	for _, r := range []rune(body) {
		buf := []byte(string(r))
		if len(digestBuf) == length {
			return string(digestBuf)
		}
		if len(buf) > 1 {
			digestBuf = append(digestBuf, r)
		}
	}
	return string(digestBuf) + misc.StrRandomCut(body, length-len(digestBuf))
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

func outputHandler(URL *url.URL, keyword string, m map[string]string) {
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
		if key == "FingerPrint" || key == "FoundDomain" {
			continue
		}
		m[key] = misc.StrRandomCut(value, 24)
	}
	fingerPrint := color.StrMapRandomColor(m, true, importantKey, varyImportantKey)
	fingerPrint = misc.FixLine(fingerPrint)
	format := "%-30v %-" + strconv.Itoa(misc.AutoWidth(color.Clear(keyword), 26+color.Count(keyword))) + "v %s"
	printStr := fmt.Sprintf(format, URL.String(), keyword, fingerPrint)
	slog.Println(slog.DATA, printStr)

	if jw := app.Setting.OutputJson; jw != nil {
		sourceMap["URL"] = URL.String()
		sourceMap["Keyword"] = keyword
		jw.Push(sourceMap)
	}
	if cw := app.Setting.OutputCSV; cw != nil {
		sourceMap["URL"] = URL.String()
		sourceMap["Keyword"] = keyword
		delete(sourceMap, "Header")
		delete(sourceMap, "Cert")
		delete(sourceMap, "Response")
		delete(sourceMap, "Body")
		sourceMap["Digest"] = strconv.Quote(sourceMap["Digest"])
		for key, value := range sourceMap {
			sourceMap[key] = chinese.ToUTF8(value)
		}
		cw.Push(sourceMap)
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
	wg.Done()
}

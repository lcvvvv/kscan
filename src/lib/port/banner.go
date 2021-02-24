package port

import (
	"app/config"
	"app/finger"
	"app/params"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"lib/iconhash"
	"lib/misc"
	"lib/shttp"
	"lib/slog"
	"lib/urlparse"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

type portInfo struct {
	Url           string
	Netloc        string
	Portid        int
	Protocol      string
	Banner        string
	Title         string
	HeaderInfo    string
	HashFinger    value
	KeywordFinger value
	Info          string
	Alive         bool
}

type value struct {
	result   string
	errorMsg error
}

func GetBanner(s string) []portInfo {
	var portInfoArr []portInfo
	Url, _ := urlparse.Load(s)
	if Url.Scheme != "" {
		for _, path := range params.SerParams.Path {
			portInfoArr = append(portInfoArr, getUrlBanner(s+"/"+path))
		}
		return portInfoArr
	}
	result := getTcpBanner(s)
	if result.Alive {
		if misc.IsInIntArr(config.Config.SslPorts, misc.Str2Int(Url.Port)) {
			//如果在常见https端口清单里面则尝试进行HTTPS协议
			s = fmt.Sprintf("https://%s", s)
			portInfoArr = append(portInfoArr, getUrlBanner(s))
		} else {
			_, isExist := config.Config.UnWebPorts[Url.Port]
			if !isExist {
				//如果不在其他常见协议的常见端口里面则尝试http协议
				s = fmt.Sprintf("http://%s", s)
				rResult := getUrlBanner(s)
				if rResult.Alive {
					//如果扫描结果正确，则是http协议端口
					portInfoArr = append(portInfoArr, rResult)
				} else {
					//如果不正确则保持原有TCP扫描结果
					portInfoArr = append(portInfoArr, result)
				}
			} else {
				//如果在其他常见协议端口里面保持原有TCP扫描结果
				portInfoArr = append(portInfoArr, result)
			}
			//if strings.Contains(result.Banner, "HTTP") {
			//	s = fmt.Sprintf("http://%s", s)
			//	portInfoArr = append(portInfoArr, getUrlBanner(s))
			//} else {
			//	portInfoArr = append(portInfoArr, result)
			//}
		}
	}
	for _, PortInfo := range portInfoArr {
		PortInfo.Info = makeResultInfo(PortInfo)
		slog.Infoln(PortInfo.Info)
	}
	return portInfoArr
}

func getUrlBanner(s string) portInfo {
	var res portInfo
	url, _ := urlparse.Load(s)
	res.Url = s
	res.Netloc = url.Host
	res.Portid = misc.Str2Int(url.Port)
	res.Protocol = getProtocol(s)
	resp, err := shttp.Get(s)
	if err != nil {
		res.Alive = false
		if strings.Contains(err.Error(), "too many") {
			//发现存在线程过高错误
			slog.Errorf("当前线程过高，请降低线程!或者请执行\"ulimit -n 50000\"命令放开操作系统限制,MAC系统可能还需要执行：\"launchctl limit maxfiles 50000 50000\"")
		}
		if strings.Contains(err.Error(), "server gave HTTP response") {
			//HTTP协议重新获取指纹
			return getUrlBanner(fmt.Sprintf("http://%s:%s", url.Host, url.Port))
		}
		if strings.Contains(err.Error(), "malformed HTTP response") {
			//TCP协议重新获取banner
			return getTcpBanner(fmt.Sprintf("%s:%s", url.Host, url.Port))
		}
		slog.Debugln(err.Error())
		return res
	}
	res.Alive = true
	query, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		slog.Debugln(err.Error())
		res.Alive = false
		return res
	}
	res.Title = getTitle(query)
	res.Banner = getHttpBanner(query)
	res.HeaderInfo = getHeaderinfo(resp.Header.Clone())
	res.HashFinger = getFingerByHash(s)
	res.KeywordFinger = getFingerByKeyword(resp)
	//res.Info = makeResultInfo(res)
	return res
}

func makeResultInfo(res portInfo) string {
	if !res.Alive {
		return ""
	}
	var infoArr []string
	if res.HashFinger.errorMsg == nil {
		infoArr = append(infoArr, "icon:"+res.HashFinger.result)
	}
	if res.KeywordFinger.errorMsg == nil {
		infoArr = append(infoArr, "keyword:"+res.KeywordFinger.result)
	}
	if res.Protocol != "" {
		infoArr = append(infoArr, res.Protocol)
	}
	if res.Portid != 0 {
		infoArr = append(infoArr, misc.Int2Str(res.Portid))
	}
	if res.HeaderInfo != "" {
		infoArr = append(infoArr, res.HeaderInfo)
	}
	Banner := ""
	if len(res.Banner) > 30 {
		i := rand.Intn(len(res.Banner) - 30)
		Banner = res.Banner[i : i+30]
	}
	res.Info = fmt.Sprintf("\r[+]%s\t%s\t%s\t%s\n", res.Url, res.Title, Banner, strings.Join(infoArr, ","))
	if params.OutPutFile != nil {
		_, _ = params.OutPutFile.WriteString(fmt.Sprintf("\r[+]%s\t%s\t%s\t%s", res.Url, res.Title, Banner, strings.Join(infoArr, ",")))
	}
	return res.Info
}

func getTcpBanner(s string) portInfo {
	var res portInfo
	url, _ := urlparse.Load(s)
	res.Url = s
	res.Netloc = url.Host
	res.Portid = misc.Str2Int(url.Port)
	res.Protocol = getProtocol(s)
	conn, err := net.DialTimeout("tcp", s, time.Second*time.Duration(params.SerParams.Timeout))
	if err != nil {
		res.Alive = false
		res.Banner = ""
		if strings.Contains(err.Error(), "too many") {
			//发现存在线程过高错误
			slog.Errorf("当前线程过高，请降低线程!或者请执行\"ulimit -n 50000\"命令放开操作系统限制,MAC系统可能还需要执行：\"launchctl limit maxfiles 50000 50000\"")
		}
		slog.Debugln(err.Error())
	} else {
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(params.SerParams.Timeout)))
		res.Alive = true
		res.KeywordFinger.errorMsg = errors.New("非Web端口")
		res.HashFinger.errorMsg = errors.New("非Web端口")
		_, _ = conn.Write([]byte("test\r\n"))
		Bytes := make([]byte, 1024)
		i, _ := conn.Read(Bytes)
		res.Banner = string(Bytes[:i])
		res.Banner = misc.FixLine(res.Banner)
		conn.Close()
	}
	return res
}

func getProtocol(s string) string {
	url, _ := urlparse.Load(s)
	if url.Scheme != "" {
		return url.Scheme
	}
	if misc.IsInIntArr(config.Config.SslPorts, misc.Str2Int(url.Port)) {
		return "https"
	}
	_, isExist := config.Config.UnWebPorts[url.Port]
	if isExist {
		return config.Config.UnWebPorts[url.Port]
	}
	return "unknow"
}

func getTitle(query *goquery.Document) string {
	result := query.Find("title").Text()
	result = misc.FixLine(result)
	//Body.Close()
	return result
}

func getHeaderinfo(header http.Header) string {
	if header.Get("SERVER") != "" {
		return header.Get("SERVER")
	}
	return ""
}

func getHttpBanner(query *goquery.Document) string {
	query.Find("script").Each(func(_ int, tag *goquery.Selection) {
		tag.Remove() // 把无用的 tag 去掉
	})
	query.Find("style").Each(func(_ int, tag *goquery.Selection) {
		tag.Remove() // 把无用的 tag 去掉
	})
	query.Find("textarea").Each(func(_ int, tag *goquery.Selection) {
		tag.Remove() // 把无用的 tag 去掉
	})
	var result string
	query.Each(func(_ int, tag *goquery.Selection) {
		result = result + tag.Text()
	})
	result = misc.FixLine(result)
	return result
}

func getFingerByKeyword(resp *http.Response) value {
	var result value
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		result.errorMsg = err
		return result
	}
	bodyStr := string(body)
	headerStr := shttp.Header2String(resp.Header)
	for _, keywordFinger := range finger.KeywordFingers.KeywordFingers {
		if keywordFinger.Type == "body" {
			if strings.Contains(bodyStr, keywordFinger.Keyword) {
				result.result = keywordFinger.Cms
				return result
			}
		}
		if keywordFinger.Type == "header" {
			if strings.Contains(headerStr, keywordFinger.Keyword) {
				result.result = keywordFinger.Cms
				return result
			}
		}
	}
	result.errorMsg = errors.New("关键字指纹库中无该指纹")
	return result
}

func getFingerByHash(s string) value {
	var iconUrl string
	var result value
	iconUrlArr, _ := urlparse.Load(s)
	if iconUrlArr.Port != "" {
		iconUrl = fmt.Sprintf("%s://%s/favicon.ico", iconUrlArr.Scheme, iconUrlArr.Host)
	} else {
		iconUrl = fmt.Sprintf("%s://%s:%s/favicon.ico", iconUrlArr.Scheme, iconUrlArr.Host, iconUrlArr.Port)
	}
	resp, err := shttp.Get(iconUrl)
	if err != nil {
		result.errorMsg = err
		return result
	}
	if resp.StatusCode != 200 {
		_ = resp.Body.Close()
		result.errorMsg = errors.New("该网站没有图标文件")
		return result
	}
	hash, _ := iconhash.Get(resp.Body)
	for _, hashFinger := range finger.HashwordFingers.HashFingers {
		if hash == hashFinger.Hash {
			result.result = hashFinger.Cms
			break
		}
	}
	if result.result == "" {
		result.errorMsg = errors.New("数据库中无该网站图标指纹")
		return result
	}
	_ = resp.Body.Close()
	return result
}

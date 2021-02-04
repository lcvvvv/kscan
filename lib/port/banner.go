package port

import (
	"../../app/config"
	"../../app/finger"
	"../../app/params"
	"../iconhash"
	"../misc"
	"../shttp"
	"../urlparse"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
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
			s = fmt.Sprintf("https://%s", s)
			portInfoArr = append(portInfoArr, getUrlBanner(s))
		} else {
			if strings.Contains(result.Banner, "HTTP") {
				s = fmt.Sprintf("http://%s", s)
				portInfoArr = append(portInfoArr, getUrlBanner(s))
			}
		}
	}
	for _, PortInfo := range portInfoArr {
		fmt.Print("\r", strings.Repeat(" ", 70))
		fmt.Print(PortInfo.Info)
	}
	return append(portInfoArr, result)
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
		fmt.Print("\r[-]", err, "\n")
		res.Alive = false
		return res
	}
	res.Alive = true
	query, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		//fmt.Print(err, "\n")
		res.Alive = false
		return res
	}
	res.Title = getTitle(query)
	res.Banner = getHttpBanner(query)
	res.HeaderInfo = getHeaderinfo(resp.Header.Clone())
	res.HashFinger = getFingerByHash(s)
	res.KeywordFinger = getFingerByKeyword(resp)
	res.Info = makeResultInfo(res)
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
	return res.Info
}

func getTcpBanner(s string) portInfo {
	var res portInfo
	url, _ := urlparse.Load(s)
	res.Url = s
	res.Netloc = url.Host
	res.Portid = misc.Str2Int(url.Port)
	res.Protocol = getProtocol(s)
	conn, err := net.DialTimeout("tcp", s, time.Second*2)
	if err != nil {
		res.Alive = false
		res.Banner = ""
	} else {
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		res.Alive = true
		_, _ = conn.Write([]byte("test\r\n"))
		Bytes := make([]byte, 1024)
		i, _ := conn.Read(Bytes)
		res.Banner = string(Bytes[:i])
		res.Banner = strings.Replace(res.Banner, "\n", "", -1)
		conn.Close()
	}
	res.Info = makeResultInfo(res)
	//fmt.Printf("[+]%s\t%d\t%s\t%s\t%s\t%s\t%s\n", s, res.portid, res.protocol, res.title, res.banner, res.hashfinger, res.keywordfinger)
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

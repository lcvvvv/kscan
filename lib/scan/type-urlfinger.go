package scan

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/lcvvvv/urlparse"
	"io"
	"io/ioutil"
	"kscan/lib/httpfinger"
	"kscan/lib/iconhash"
	"kscan/lib/misc"
	"kscan/lib/shttp"
	"kscan/lib/slog"
	"math/rand"
	"net/http"
)

type HttpFinger struct {
	StatusCode     int
	Response       string
	ResponseDigest string
	Title          string
	Header         string
	HeaderInfo     string
	HashFinger     string
	KeywordFinger  string
	Info           string
}

func NewHttpFinger() *HttpFinger {
	return &HttpFinger{
		StatusCode:     0,
		Response:       "",
		ResponseDigest: "",
		Title:          "",
		Header:         "",
		HashFinger:     "",
		KeywordFinger:  "",
		Info:           "",
	}
}

func (h *HttpFinger) LoadHttpResponse(url *urlparse.URL, resp *http.Response) {
	h.Title = getTitle(shttp.GetBody(resp))
	h.Header = getHeader(resp.Header.Clone())
	h.HeaderInfo = getHeaderinfo(resp.Header.Clone())
	h.Response = getResponse(shttp.GetBody(resp))
	h.ResponseDigest = getResponseDigest(shttp.GetBody(resp))
	h.HashFinger = getFingerByHash(*url)
	h.KeywordFinger = getFingerByKeyword(h.Header, h.Response)
	h.makeInfo()
	_ = resp.Body.Close()
}

func getTitle(resp io.Reader) string {
	query, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		slog.Debug(err.Error())
	}
	result := query.Find("title").Text()
	result = misc.FixLine(result)
	//Body.Close()
	return result
}

func getHeader(header http.Header) string {
	return shttp.Header2String(header)
}

func getResponse(resp io.Reader) string {
	body, err := ioutil.ReadAll(resp)
	if err != nil {
		slog.Debug(err.Error())
		return ""
	}
	bodyStr := string(body)
	return bodyStr
}

func getResponseDigest(resp io.Reader) string {
	var result string
	query, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		slog.Debug(err.Error())
	}
	query.Find("script").Each(func(_ int, tag *goquery.Selection) {
		tag.Remove() // 把无用的 tag 去掉
	})
	query.Find("style").Each(func(_ int, tag *goquery.Selection) {
		tag.Remove() // 把无用的 tag 去掉
	})
	query.Find("textarea").Each(func(_ int, tag *goquery.Selection) {
		tag.Remove() // 把无用的 tag 去掉
	})
	query.Each(func(_ int, tag *goquery.Selection) {
		result = result + tag.Text()
	})
	result = misc.FixLine(result)
	if len(result) > 10 {
		return result
	} else {
		bodyBuf, _ := ioutil.ReadAll(resp)
		return misc.FixLine(string(bodyBuf))
	}
}

func getHeaderinfo(header http.Header) string {
	if header.Get("SERVER") != "" {
		return "server:" + header.Get("SERVER")
	}
	return ""
}
func getFingerByKeyword(header string, body string) string {
	return httpfinger.KeywordFinger.Match(body, header)
}

func getFingerByHash(url urlparse.URL) string {
	url.Path = "/favicon.ico"
	resp, err := shttp.Get(url.UnParse())
	if err != nil {
		slog.Debug(url.UnParse() + err.Error())
		return ""
	}
	if resp.StatusCode != 200 {
		slog.Debug(url.UnParse() + "没有图标文件")
		return ""
	}
	hash, err := iconhash.Get(resp.Body)
	if err != nil {
		slog.Debug(url.UnParse() + err.Error())
		return ""
	}
	_ = resp.Body.Close()
	return httpfinger.FaviconHash.Match(hash)
}

func (h *HttpFinger) makeInfo() {
	var info string
	if h.HashFinger != "" {
		info += ",favicon:" + h.HashFinger
	}
	if h.KeywordFinger != "" {
		info += ",keyword:" + h.KeywordFinger
	}
	if h.HeaderInfo != "" {
		info += "," + h.HeaderInfo
	}
	if h.ResponseDigest != "" {
		if len(h.ResponseDigest) > 30 {
			i := rand.Intn(len(h.ResponseDigest) - 30)
			fmt.Println(i)
			info += "," + h.ResponseDigest[i:i+30]
		} else {
			info += "," + h.ResponseDigest
		}
	}
	if info != "" {
		h.Info = info[1:]
	}
}

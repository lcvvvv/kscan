package shttp

import (
	"app/config"
	"app/params"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"lib/misc"
	"lib/slog"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var headerKeys = map[string]string{
	"Server":                           "中间件名称",
	"X-Powered-By":                     "中间件或开发语言名称",
	"Content-Length":                   "返回包长度",
	"Last-Modified":                    "最后一次验证日期",
	"Etag":                             "连接标签",
	"Accept-Ranges":                    "不知道",
	"Sdn-Via":                          "不知道",
	"Use_new_portal":                   "不知道",
	"X-Robots-Tag":                     "不知道",
	"X-Application-Context":            "不知道",
	"X-Content-Type-Options":           "不知道",
	"X-Xss-Protection":                 "不知道",
	"Date":                             "日期",
	"Expires":                          "失效日期",
	"Content-Type":                     "正文类型",
	"Set-Cookie":                       "设置cookie值",
	"Connection":                       "连接类型",
	"Vary":                             "不知道",
	"Keep-Alive":                       "长链接保存时间",
	"X-Frame-Options":                  "框架选项",
	"X-Aspnet-Version":                 "keyword:ASP.NET",
	"X-Aspnetmvc-Version":              "keyword:ASP.NET MVC",
	"Content-Language":                 "正文语言",
	"Cache-Control":                    "缓存控制",
	"Pragma":                           "程序类型",
	"Progma":                           "不知道",
	"Access-Control-Allow-Origin":      "同源策略",
	"Access-Control-Allow-Methods":     "同源策略",
	"Access-Control-Allow-Headers":     "同源策略",
	"Access-Control-Expose-Headers":    "同源策略",
	"Access-Control-Allow-Credentials": "同源策略",
	"X-Enterprise":                     "不知道",
	"X-Lang":                           "不知道",
	"X-Timezone":                       "时区",
	"X-Arch":                           "系统架构",
	"X-Support-Wifi":                   "不知道",
	"X-Timestamp":                      "时间戳",
	"X-Sysbit":                         "不知道",
	"X-Support-I18n":                   "不知道",
	"P3p":                              "不知道",
	"Content-Security-Policy":          "不知道",
	"Www-Authenticate":                 "认证参数",
	"X-Ua-Compatible":                  "不知道",
	"Entry1":                           "不知道",
	"Accept-Encoding":                  "不知道",
	"X-Amz-Request-Id":                 "不知道",
	"Cfl_asynch":                       "keyword:浙江大华dh650平台设备",
}

var newHeaderKeys *os.File

func initNewHeaderKeys() *os.File {
	if newHeaderKeys != nil {
		return newHeaderKeys
	}
	return misc.SafeOpen("newHeaderkeys.txt")
}

func Get(Url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		slog.Debug(err.Error())
		return nil, err
	}
	request.Header.Add("User-Agent", getUserAgent())
	request.Header.Add("Cookie", "rememberMe=b69375edcb2b3c5084c02bd9690b6625")
	request.Close = true
	//修改Host头部参数
	if params.SerParams.Host != "" {
		request.Header.Add("Host", params.SerParams.Host)
	}
	tr := &http.Transport{}
	if params.SerParams.Proxy != "" {
		uri, _ := url.Parse(params.SerParams.Proxy)
		(*tr).Proxy = http.ProxyURL(uri)
	}
	(*tr).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
	}
	(*tr).DisableKeepAlives = false
	client := &http.Client{}
	client.Transport = tr
	client.Timeout = time.Second * time.Duration(params.SerParams.Timeout)
	resp, err := client.Do(request)
	if err == nil {
		//校验http头部
		newHeaderKeys = initNewHeaderKeys()
		for key := range resp.Header {
			if headerKeys[key] == "" {
				headerKeys[key] = "New"
				_, _ = newHeaderKeys.WriteString(fmt.Sprintf("%s: %s\n", key, resp.Header.Get(key)))
				slog.Warningf("\r[*]发现生僻Http头部：%s: %s", key, resp.Header.Get(key))
			}
		}
		//校验http状态码
		if len(config.Config.HttpCode) > 0 {
			if !misc.IsInIntArr(config.Config.HttpCode, resp.StatusCode) {
				resp = nil
				err = errors.New("HttpStatusCode不在范围内")
				return resp, err
			}
		}
		//修复乱码问题
		body2UTF8(resp)

		//检测内容是否压缩
		//if resp.Header.Get("Content-Encoding") == "gzip" {
		//	var reader io.ReadCloser
		//	reader, err = gzip.NewReader(resp.Body)
		//	if err != nil {
		//		resp.Body = reader
		//	}
		//}
	}
	return resp, err
}

func getUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(config.Config.UserAgents))
	return config.Config.UserAgents[i]
}

func getCharset(contentType string) string {
	if len(contentType) > 0 {
		charsetRegexp, _ := regexp.Compile("charset=(.+)[;$]?")
		charset := charsetRegexp.FindStringSubmatch(contentType)
		if len(charset) == 2 {
			return strings.ToLower(charset[1])
		} else {
			return "unknown"
		}
	} else {
		return "unknown"
	}
}

func Header2String(header http.Header) string {
	var result string
	for i := range header {
		result = strings.Join([]string{result, fmt.Sprintf("%s: %s\n", i, header.Get(i))}, "")
	}
	return result
}

func body2UTF8(resp *http.Response) {
	charset := getCharset(resp.Header.Get("Content-Type"))
	bodyBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		slog.Debug(err.Error())
	}
	if strings.Contains(charset, "gb") {
		utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(bodyBuf)
		resp.Body = ioutil.NopCloser(bytes.NewReader(utf8Data))
		return
	}
	if charset == "unknown" {
		if isUtf8(bodyBuf) {
			resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBuf))
		} else {
			utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(bodyBuf)
			resp.Body = ioutil.NopCloser(bytes.NewReader(utf8Data))
		}
		return
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBuf))
	return
}

func GetBody(resp *http.Response) io.Reader {
	bodyBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		slog.Debug(err.Error())
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBuf))
	return bytes.NewReader(bodyBuf)
}

func isUtf8(data []byte) bool {
	for i := 0; i < len(data); {
		if data[i]&0x80 == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if data[i]&0xc0 != 0x80 {
					return false
				}
				i++
			}
		} else {
			//其他情况说明不是utf-8
			return false
		}
	}
	return true
}

func preNUm(data byte) int {
	str := fmt.Sprintf("%b", data)
	var i int = 0
	for i < len(str) {
		if str[i] != '1' {
			break
		}
		i++
	}
	return i
}

package httpfinger

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/PuerkitoBio/goquery"
	"github.com/lcvvvv/appfinger/iconhash"
	"github.com/lcvvvv/simplehttp"
	"github.com/lcvvvv/stdio/chinese"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Banner struct {
	Protocol string
	Port     string
	Header   string
	Body     string
	Response string
	Cert     string
	Title    string
	Hash     string
	Icon     string
	ICP      string
}

func GetBannerWithURL(URL *url.URL, req *http.Request, cli *http.Client) (*Banner, error) {
	var banner = &Banner{}
	if req == nil {
		req, _ = simplehttp.NewRequest(http.MethodGet, URL.String(), nil)
	}

	req.Header.Set("User-Agent", simplehttp.RandomUserAgent())

	if cli == nil {
		cli = simplehttp.NewClient()
	}

	resp, err := simplehttp.Do(cli, req)
	if err != nil {
		return nil, err
	}

	var body = chinese.ToUTF8(resp.Raw.Body)

	banner.Protocol = URL.Scheme
	banner.Port = GetURLPortString(URL)
	banner.Header = resp.Raw.Header
	banner.Body = body
	banner.Response = resp.Raw.String()
	banner.Cert = resp.Raw.Cert
	banner.Title = getTitle(body)
	banner.Hash = getHash(body)
	banner.Icon = getIcon(*URL, body)
	banner.ICP = getICP(body)

	return banner, err
}

func GetBannerWithResponse(URL *url.URL, response string, req *http.Request, cli *http.Client) (*Banner, error) {
	if URL.Scheme == "https" {
		return GetBannerWithURL(URL, req, cli)
	}
	if statusCode := getStatusCode(response); statusCode >= 300 && statusCode <= 400 {
		return GetBannerWithURL(URL, req, cli)
	}
	header, body := simplehttp.SplitHeaderAndBody(response)
	body = chinese.ToUTF8(body)

	var banner = &Banner{}
	banner.Protocol = URL.Scheme
	banner.Port = GetURLPortString(URL)
	banner.Header = header
	banner.Body = body
	banner.Response = response
	banner.Cert = ""
	banner.Title = getTitle(body)
	banner.Hash = getHash(body)
	banner.Icon = getIcon(*URL, body)
	banner.ICP = getICP(body)
	return banner, nil
}

var (
	regxIconPath           = regexp.MustCompile(`^/.*$`)
	regxIconURL            = regexp.MustCompile(`^(?:http|https)://.*$`)
	regxIconURLNotProtocol = regexp.MustCompile(`^://.*$`)
)

func getIcon(URL url.URL, body string) string {
	path := getIconPath(body)
	if regxIconPath.MatchString(path) == true {
		URL.Path = ""
		return getIconHash(URL.String() + path)
	}
	if regxIconURL.MatchString(path) == true {
		return getIconHash(path)
	}
	if regxIconURLNotProtocol.MatchString(path) == true {
		return getIconHash(URL.Scheme + path)
	}
	return ""
}

func getIconPath(body string) string {
	bodyReader := strings.NewReader(body)
	query, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return "/favicon.ico"
	}
	selector := query.Find("link[rel='shortcut icon']")
	if iconHref, ok := selector.Attr("href"); ok == true {
		return iconHref
	}
	return "/favicon.ico"
}

func getIconHash(URL string) string {
	resp, err := simplehttp.Get(URL)
	if err != nil {
		return ""
	}
	hash := iconhash.Encode([]byte(resp.Raw.Body))
	return hash
}

func getTitle(body string) string {
	bodyReader := strings.NewReader(body)
	query, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return ""
	}
	title := query.Find("title").Text()
	title = regexp.MustCompile(`\s|\n|\r`).ReplaceAllString(title, "")
	return title
}

func getHash(body string) string {
	hash := md5.New()
	hash.Write([]byte(body))
	return hex.EncodeToString(hash.Sum(nil))
}

var (
	provinces = []string{
		"京", "津", "冀", "晋", "蒙", "辽", "吉", "黑",
		"沪", "苏", "浙", "皖", "闽", "赣", "鲁", "豫",
		"湘", "粤", "桂", "琼", "川", "蜀", "贵", "黔",
		"云", "滇", "渝", "藏", "陕", "秦", "甘", "陇",
		"青", "宁", "新", "港", "澳", "台", "鄂",
	}
	provincesString = strings.Join(provinces, "|")
	icpRegx         = regexp.MustCompile(`(?:` + provincesString + `)ICP备\s*\d+号(?:-\d+)?`)
)

func getICP(body string) string {
	if icpRegx.MatchString(body) == true {
		return icpRegx.FindString(body)
	}
	return ""
}
func GetURLPortString(URL *url.URL) string {
	if URL.Port() != "" {
		port, _ := strconv.Atoi(URL.Port())
		return strconv.Itoa(port)
	}
	if URL.Scheme == "https" {
		return "443"
	}
	if URL.Scheme == "http" {
		return "80"
	}
	return ""
}

var statusCodeRegx = regexp.MustCompile(`^HTTP/\d\.\d (\d\d\d).*`)

func getStatusCode(responseRaw string) int {
	if statusCodeRegx.MatchString(responseRaw) == false {
		return 0
	}
	statusCodeString := statusCodeRegx.FindAllStringSubmatch(responseRaw, -1)[0][1]
	statusCode, _ := strconv.Atoi(statusCodeString)
	return statusCode
}

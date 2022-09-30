package appfinger

import (
	"crypto/md5"
	"encoding/hex"
	"net"
	"reflect"
	"regexp"
	"strings"
)

type Banner struct {
	Header   string
	Body     string
	Response string
	Cert     string
	Title    string
	Hash     string
	Icon     string
	ICP      string

	FoundDomain string
	FoundIP     string
}

func convBanner(other interface{}) *Banner {
	if reflect.ValueOf(other).IsNil() == true {
		return nil
	}

	var banner = &Banner{}
	bannerT := reflect.TypeOf(banner)
	bannerV := reflect.ValueOf(banner)

	t := reflect.TypeOf(other)
	v := reflect.ValueOf(other)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < bannerT.Elem().NumField(); i++ {
		var key = bannerT.Elem().Field(i).Name
		if t, ok := t.FieldByName(key); ok == false || t.Type.Kind() != reflect.String {
			continue
		}
		value := v.FieldByName(key).String()
		bannerValue := bannerV.Elem().FieldByName(key)
		bannerValue.SetString(value)
	}

	banner.FoundIP = findIP(banner.Body, banner.Cert)
	banner.FoundDomain = findDomain(banner.Body, banner.Cert)
	return banner
}

func convBannerWithRaw(response string) *Banner {
	return &Banner{
		Header:   "",
		Body:     "",
		Response: response,
		Cert:     "",
		Title:    "",
		Hash:     getHash(response),
		Icon:     "",
	}
}

func getHash(body string) string {
	hash := md5.New()
	hash.Write([]byte(body))
	return hex.EncodeToString(hash.Sum(nil))
}

var (
	domainRoot = []string{
		"com", "net", "org", "aero", "biz", "coop", "info", "museum", "name", "pro", "top", "xyz",
		"loan", "wang", "ad", "ae", "af", "ag", "ai", "al", "am", "an", "ao", "aq", "ar", "as", "at",
		"au", "aw", "az", "ba", "bb", "bd", "be", "bf", "bg", "bh", "bi", "bj", "bm", "bn", "bo", "br",
		"bs", "bt", "bv", "bw", "by", "bz", "ca", "cc", "cf", "cg", "ch", "ci", "ck", "cl", "cm", "cn",
		"co", "cq", "cr", "cu", "cv", "cx", "cy", "cz", "de", "dj", "dk", "dm", "do", "dz", "ec", "ee",
		"eg", "eh", "es", "et", "ev", "fi", "fj", "fk", "fm", "fo", "fr", "ga", "gb", "gd", "ge", "gf",
		"gh", "gi", "gl", "gm", "gn", "gp", "gr", "gt", "gu", "gw", "gy", "hk", "hm", "hn", "hr", "ht",
		"hu", "id", "ie", "il", "in", "io", "iq", "ir", "is", "it", "jm", "jo", "jp", "ke", "kg", "kh",
		"ki", "km", "kn", "kp", "kr", "kw", "ky", "kz", "la", "lb", "lc", "li", "lk", "lr", "ls", "lt",
		"lu", "lv", "ly", "ma", "mc", "me", "md", "mg", "mh", "ml", "mm", "mn", "mo", "mp", "mq", "mr",
		"ms", "mt", "mv", "mw", "mx", "my", "mz", "na", "nc", "ne", "nf", "ng", "ni", "nl", "no", "np",
		"nr", "nt", "nu", "nz", "om", "pa", "pe", "pf", "pg", "ph", "pk", "pl", "pm", "pn", "pr", "pt",
		"pw", "py", "qa", "re", "ro", "ru", "rw", "sa", "sb", "sc", "sd", "se", "sg", "sh", "si", "sj",
		"sk", "sl", "sm", "sn", "so", "sr", "st", "su", "sy", "sz", "tc", "td", "tf", "tg", "th", "tj",
		"tk", "tl", "tm", "tn", "to", "tp", "tr", "tt", "tv", "tw", "tz", "ua", "ug", "uk", "us", "uy",
		"uz", "va", "vc", "ve", "vg", "vn", "vu", "wf", "ws", "ye", "yu", "za", "zm", "zr", "zw"}
	domainRootString = strings.Join(domainRoot, "|")
	domainRegx       = regexp.MustCompile(`([a-zA-Z0-9][-a-zA-Z0-9]{0,62}(?:\.[a-zA-Z0-9][-a-zA-Z0-9]{1,62})*\.(?:` + domainRootString + `))[^a-zA-Z0-9-]`)
)

func findDomain(strArr ...string) string {
	var domains []string
	for _, str := range strArr {
		if domainRegx.MatchString(str) == true {
			for _, value := range domainRegx.FindAllString(str, -1) {
				length := len(value)
				if length > 253 {
					continue
				}
				domains = append(domains, value[:length-1])
			}
		}
	}
	domains = removeDuplicate(domains)
	return strings.Join(domains, "、")
}

var (
	ipRegx = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
)

func findIP(strArr ...string) string {
	var IPs []string
	for _, str := range strArr {
		if ipRegx.MatchString(str) == true {
			for _, value := range ipRegx.FindAllString(str, -1) {
				if net.ParseIP(value) != nil {
					continue
				}
				IPs = append(IPs, value)
			}
		}
	}
	IPs = removeDuplicate(IPs)
	return strings.Join(IPs, ",")
}

func removeDuplicate(p []string) []string {
	result := make([]string, 0, len(p))
	temp := map[string]struct{}{}
	for _, item := range p {
		if _, ok := temp[item]; !ok { //如果字典中找不到元素，ok=false，!ok为true，就往切片中append元素。
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

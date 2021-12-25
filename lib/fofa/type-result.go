package fofa

import (
	"fmt"
	"kscan/lib/misc"
	"reflect"
	"regexp"
	"strings"
)

type Result struct {
	Host, Title, Ip, Domain, Port, Country string
	Province, City, Country_name, Protocol string
	Server, Banner, Isp, As_organization   string
	Header, Cert                           string
}

func (r *Result) Fix() {
	if r.Protocol != "" {
		r.Host = fmt.Sprintf("%s://%s:%s", r.Protocol, r.Ip, r.Port)
	}
	if regexp.MustCompile("http([s]?)://.*").MatchString(r.Host) == false && r.Protocol == "" {
		r.Host = "http://" + r.Host
	}
	if r.Banner != "" {
		r.Banner = misc.FixLine(r.Banner)
		r.Banner = misc.StrRandomCut(r.Banner, 20)
	}
	if r.Title == "" && r.Protocol != "" {
		r.Title = strings.ToUpper(r.Protocol)
	}

	r.Title = misc.FixLine(r.Title)

}

func (r Result) Map() map[string]string {
	t := reflect.TypeOf(r)
	v := reflect.ValueOf(r)
	m := make(map[string]string)
	for k := 0; k < t.NumField(); k++ {
		key := t.Field(k).Name
		value := v.Field(k).String()
		m[key] = value
	}
	return m
}

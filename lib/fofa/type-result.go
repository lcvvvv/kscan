package fofa

import (
	"reflect"
)

type Result struct {
	Host, Title, Ip, Domain, Port, Country string
	Province, City, Country_name, Protocol string
	Server, Banner, Isp, As_organization   string
	Header, Cert                           string
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

package urlparse

import (
	"errors"
	"regexp"
)

type Url struct {
	Scheme, Host, Port, Path string
}

func Load(s string) (Url, error) {
	r := regexp.MustCompile("^(?:(http|https)://)?([A-Za-z0-9.\\-]+(?:\\.[A-Za-z0-9.\\-]+))(?::(\\d+))?(/*[\\w/%]*)?$")
	o := r.FindStringSubmatch(s)
	if len(o) != 5 {
		return Url{}, errors.New("URL格式不正确")
	}
	if o[3] == "" {
		switch o[1] {
		case "https":
			o[3] = "443"
		case "http":
			o[3] = "80"
		}
	}
	return Url{
		Scheme: o[1],
		Host:   o[2],
		Port:   o[3],
		Path:   o[4],
	}, nil
}

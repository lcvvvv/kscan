package urlparse

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type URL struct {
	Scheme, Netloc, Path string
	Port                 int
}

func Load(s string) (*URL, error) {
	r := regexp.MustCompile("^(?:(http|https)://)?([A-Za-z0-9.\\-]+(?:\\.[A-Za-z0-9.\\-]+))(?::(\\d+))?(/*[\\w/%]*)?$")
	o := r.FindStringSubmatch(s)

	if len(o) != 5 {
		return nil, errors.New("URL格式不正确")
	}

	Scheme := o[1]
	Netloc := o[2]
	Port, _ := strconv.Atoi(o[3])
	Path := o[4]

	if Port == 0 {
		switch Scheme {
		case "https":
			Port = 443
		case "http":
			Port = 80
		}
	}

	return &URL{
		Scheme: Scheme,
		Netloc: Netloc,
		Port:   Port,
		Path:   Path,
	}, nil
}

func (i *URL) UnParse() string {
	return fmt.Sprintf("%s://%s:%d%s", i.Scheme, i.Netloc, i.Port, i.Path)
}

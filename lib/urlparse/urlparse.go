package urlparse

import (
	"fmt"
	"net/url"
	"strconv"
)

type URL struct {
	Scheme, Netloc, Path string
	Port                 int
	url                  *url.URL
}

func Load(s string) (*URL, error) {
	r, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	return &URL{
		Scheme: r.Scheme,
		Netloc: r.Host,
		Path:   r.Path,
		Port: func() int {
			if r.Port() != "" {
				p, _ := strconv.Atoi(r.Port())
				return p
			}
			if r.Scheme == "https" {
				return 443
			}
			if r.Scheme == "http" {
				return 80
			}
			return 80
		}(),
	}, nil
}

func (i *URL) UnParse() string {
	if i.Scheme == "https" && i.Port == 443 {
		return fmt.Sprintf("https://%s%s", i.Netloc, i.Path)
	}
	if i.Scheme == "http" && i.Port == 80 {
		return fmt.Sprintf("http://%s%s", i.Netloc, i.Path)
	}
	return fmt.Sprintf("%s://%s:%d%s", i.Scheme, i.Netloc, i.Port, i.Path)
}

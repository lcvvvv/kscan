package appfinger

import (
	"errors"
	"github.com/lcvvvv/appfinger/gorpc"
	"github.com/lcvvvv/appfinger/httpfinger"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var supportProtocols = []string{
	"http",
	"https",
	"rpc",
}

var supportProtocolRegx = regexp.MustCompile("^" + strings.Join(supportProtocols, "|") + "$")

func Search(URL *url.URL, banner *Banner) *FingerPrint {
	var products, hostnames []string
	switch URL.Scheme {
	case "http":
		products = search(convHttpBanner(URL, banner))
		return &FingerPrint{products, "", "", ""}
	case "https":
		products := search(convHttpBanner(URL, banner))
		return &FingerPrint{products, "", "", ""}
	case "rpc":
		hostnames, _ = gorpc.GetHostname(URL.Hostname())
		return &FingerPrint{emptyProductName, strings.Join(hostnames, ";"), "", ""}
	}
	return nil
}

func SupportCheck(protocol string) bool {
	return supportProtocolRegx.MatchString(protocol)
}

func GetBannerWithResponse(URL *url.URL, response string, req *http.Request, cli *http.Client) (*Banner, error) {
	switch URL.Scheme {
	case "http":
		httpBanner, err := httpfinger.GetBannerWithResponse(URL, response, req, cli)
		return convBanner(httpBanner), err
	case "https":
		httpBanner, err := httpfinger.GetBannerWithResponse(URL, response, req, cli)
		return convBanner(httpBanner), err
	default:
		return convBannerWithRaw(response), nil
	}
}

func GetBannerWithURL(URL *url.URL, req *http.Request, cli *http.Client) (*Banner, error) {
	switch URL.Scheme {
	case "http":
		httpBanner, err := httpfinger.GetBannerWithURL(URL, req, cli)
		return convBanner(httpBanner), err
	case "https":
		httpBanner, err := httpfinger.GetBannerWithURL(URL, req, cli)
		return convBanner(httpBanner), err
	}
	return nil, errors.New("unsupported protocol")
}

func convHttpBanner(URL *url.URL, banner *Banner) *httpfinger.Banner {
	return &httpfinger.Banner{
		Protocol: URL.Scheme,
		Port:     URL.Port(),
		Header:   banner.Header,
		Body:     banner.Body,
		Response: banner.Response,
		Cert:     banner.Cert,
		Title:    banner.Title,
		Hash:     banner.Hash,
		Icon:     banner.Icon,
	}
}

package scanner

import (
	"errors"
	"github.com/lcvvvv/appfinger"
	"github.com/lcvvvv/gonmap"
	"net/http"
	"net/url"
)

type foo2 struct {
	URL      *url.URL
	response *gonmap.Response
	req      *http.Request
	client   *http.Client
}

const (
	NotSupportProtocol = "protocol is not support"
)

type URLClient struct {
	*client
	HandlerMatched func(url *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint)
	HandlerError   func(url *url.URL, err error)
}

func NewURLScanner(config *Config) *URLClient {
	var client = &URLClient{
		client:         newConfig(config, config.Threads),
		HandlerMatched: func(url *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint) {},
		HandlerError:   func(url *url.URL, err error) {},
	}
	client.pool.Interval = config.Interval
	client.pool.Function = func(in interface{}) {
		value := in.(foo2)
		URL := value.URL
		response := value.response
		req := value.req
		cli := value.client
		if appfinger.SupportCheck(URL.Scheme) == false {
			client.HandlerError(URL, errors.New(NotSupportProtocol))
			return
		}
		var banner *appfinger.Banner
		var finger *appfinger.FingerPrint
		var err error
		if response == nil || req != nil || cli != nil {
			banner, err = appfinger.GetBannerWithURL(URL, req, cli)
			if err != nil {
				client.HandlerError(URL, err)
				return
			}
			finger = appfinger.Search(URL, banner)
		} else {
			//banner, err = appfinger.GetBannerWithResponse(URL, response.Raw, req, cli)
			banner, err = appfinger.GetBannerWithURL(URL, req, cli)
			if err != nil {
				client.HandlerError(URL, err)
				return
			}
			finger = appfinger.Search(URL, banner)
			appendTcpBannerInFinger(finger, response.FingerPrint)
		}
		client.HandlerMatched(URL, banner, finger)
	}
	return client
}

func (c *URLClient) Push(URL *url.URL, response *gonmap.Response, req *http.Request, client *http.Client) {
	c.pool.Push(foo2{URL, response, req, client})
}

func appendTcpBannerInFinger(finger *appfinger.FingerPrint, gonmapFinger *gonmap.FingerPrint) *appfinger.FingerPrint {
	if gonmapFinger.ProductName != "" {
		if gonmapFinger.Version != "" {
			finger.AddProduct(gonmapFinger.ProductName + "/" + gonmapFinger.Version)
		}
		finger.AddProduct(gonmapFinger.ProductName)
	}

	if gonmapFinger.OperatingSystem != "" {
		finger.AddProduct(gonmapFinger.OperatingSystem)
	}

	if gonmapFinger.DeviceType != "" {
		finger.AddProduct(gonmapFinger.DeviceType)
	}

	if gonmapFinger.Info != "" {
		finger.AddProduct(gonmapFinger.Info)
	}

	if gonmapFinger.Hostname != "" {
		finger.Hostname = gonmapFinger.Hostname
	}
	return finger
}

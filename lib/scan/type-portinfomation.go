package scan

import (
	"fmt"
	"github.com/lcvvvv/gonmap"
	"github.com/lcvvvv/urlparse"
	"kscan/lib/misc"
)

type PortInformation struct {
	Target         *urlparse.URL
	Response       string
	ResponseDigest string
	Status         string
	Finger         *gonmap.Finger
	HttpFinger     *HttpFinger
	Info           string
}

func NewPortInformation(u *urlparse.URL) *PortInformation {
	return &PortInformation{
		Target:     u,
		Response:   "",
		Status:     "UNKNOWN",
		Finger:     nil,
		HttpFinger: nil,
		Info:       "",
	}
}

func (p *PortInformation) LoadGonmapPortInformation(g *gonmap.PortInfomation) {
	p.Response = g.Response()
	p.ResponseDigest = misc.MustLength(misc.FilterPrintStr(p.Response), 0)
	p.Status = g.Status()
	p.Finger = g.Finger()
	p.Target.Scheme = p.Finger.Service
}

func (p *PortInformation) LoadHttpFinger(h *HttpFinger) {
	p.HttpFinger = h
}

func (p *PortInformation) MakeInfo() {
	p.Info = "%s %d %s %s"
	target := p.Target.UnParse()
	code := len(p.Response)
	digest := p.ResponseDigest
	fingerprint := p.Finger.Info
	if p.HttpFinger != nil {
		if p.HttpFinger.StatusCode != 0 {
			fmt.Println()
			fmt.Println(p.HttpFinger.StatusCode)
			fmt.Println()
			code = p.HttpFinger.StatusCode
			if p.HttpFinger.Title != "" {
				digest = p.HttpFinger.Title
			}
			fingerprint += p.HttpFinger.Info
		}
	}
	p.Info = fmt.Sprintf(p.Info, target, code, digest, fingerprint)
}

package run

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
	ErrorMsg       error
}

func NewPortInformation(u *urlparse.URL) *PortInformation {
	return &PortInformation{
		Target:     u,
		Response:   "",
		Status:     "UNKNOWN",
		Finger:     nil,
		HttpFinger: nil,
		Info:       "",
		ErrorMsg:   nil,
	}
}

func (p *PortInformation) LoadGonmapPortInformation(g *gonmap.PortInfomation) {
	p.Status = g.Status()
	p.Finger = g.Finger()
	p.Target.Scheme = p.Finger.Service
	p.ErrorMsg = g.ErrorMsg
	if g.Service() != "CLOSED" {
		p.Response = g.Response()
		p.ResponseDigest = misc.MustLength(misc.FilterPrintStr(p.Response), 10)
	}
}

func (p *PortInformation) LoadHttpFinger(h *HttpFinger) {
	p.HttpFinger = h
}

func (p *PortInformation) MakeInfo() {
	if p.Target.Scheme == "" {
		p.Target.Scheme = "unknown"
	}
	target := p.Target.UnParse()
	code := len(p.Response)
	digest := p.ResponseDigest
	fingerprint := p.Finger.Information()
	if p.HttpFinger != nil {
		if p.HttpFinger.StatusCode != 0 {
			code = p.HttpFinger.StatusCode
			if p.HttpFinger.Title != "" {
				digest = p.HttpFinger.Title
			}
			fingerprint += "," + p.HttpFinger.Info
		}
	}
	if fingerprint != "" {
		if fingerprint[:1] == "," {
			fingerprint = fingerprint[1:]
		}
	}
	p.Info = misc.FixLine(fingerprint)
	p.Info = fmt.Sprintf("%s\t%d\t%s\t%s", target, code, digest, fingerprint)
}

func (p *PortInformation) CLOSED() *PortInformation {
	p.Status = "CLOSED"
	return p
}

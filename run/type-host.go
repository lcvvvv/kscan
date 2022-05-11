package run

import (
	"fmt"
	"github.com/lcvvvv/gonmap"
	"kscan/lib/color"
	"kscan/lib/smap"
)

type Host struct {
	icmp     bool
	tcp      bool
	portScan bool

	PortNum int

	TcpBannerCount        int
	TcpBannerUnknownCount int
	TcpBannerScanned      bool

	TcpCount     int
	OpenTcpCount int

	Map struct {
		//Map[string]*gonmap.TcpBanner
		TcpBanner *smap.SMap
		//Map[string]*gonmap.AppBanner
		//App *smap.SMap
	}

	addr string
}

func NewHost(addr string, length int) *Host {
	host := &Host{}
	host.addr = addr

	host.icmp = false

	host.tcp = false
	host.PortNum = length

	host.Map.TcpBanner = smap.New()
	//host.Map.App = smap.New()

	return host
}

func (h *Host) Up() *Host {
	h.icmp = true
	return h
}

func (h *Host) Down() *Host {
	h.icmp = false
	return h
}

func (h *Host) Alive() bool {
	if h.icmp == true {
		return true
	}
	if h.tcp == true {
		return true
	}
	return false
}

func (h *Host) AddPort(port *Port) {
	h.TcpCount++
	if port.Status == Open {
		h.OpenTcpCount++
		h.tcp = true
	}
}

func (h *Host) DisplayUnknownPorts() string {
	var portSlice []int
	var dispSlice []string

	h.Map.TcpBanner.Range(func(key, value interface{}) bool {
		tcpBanner := value.(*gonmap.TcpBanner)
		if tcpBanner.Status() == gonmap.Open || tcpBanner.Status() == gonmap.Unknown {
			portSlice = append(portSlice, tcpBanner.Target.Port())
		}
		return true
	})

	for _, port := range portSlice {
		dispSlice = append(dispSlice, fmt.Sprintf("%d[%s?]", port, gonmap.GuessProtocol(port)))
	}
	url := fmt.Sprintf("%s://%s", "unknown", h.addr)
	title := "UnknownProtocols"
	output := fmt.Sprintf("%-30v %-26v %s%s", url, title,
		"ThesePortsAreUnknownProtocols: ",
		color.StrSliceRandomColor(dispSlice))
	return output
}

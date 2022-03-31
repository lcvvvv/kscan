package run

import (
	"fmt"
	"kscan/core/gonmap"
	"kscan/lib/color"
	"kscan/lib/smap"
)

type Host struct {
	status struct {
		//主机是否存活
		alive bool
		//是否存在开放端口
		open bool
		//扫描状态
		portScan bool
		//tcpScan  bool
		//appScan  bool
	}

	Length struct {
		//需要进行端口扫描的端口数量
		Port int
		//需要进行tcp协议识别的端口数量
		Tcp int
		//需要进行应用层识别的端口数量
		//app  int
	}

	Map struct {
		//Map[string]*Port
		Port *smap.SMap
		//Map[string]*gonmap.TcpBanner
		Tcp *smap.SMap
		//Map[string]*gonmap.AppBanner
		//App *smap.SMap
	}

	addr string
}

const (
	Close   = 0x00001a
	Open    = 0x00002b
	Unknown = 0x00003c
)

func NewHost(addr string, length int) *Host {
	host := &Host{}

	host.addr = addr

	host.status.alive = false
	host.status.open = false
	host.status.portScan = false
	//host.status.tcpScan = false
	//host.status.appScan = false

	host.Length.Port = length
	host.Length.Tcp = 0
	//host.Length.App = 0

	host.Map.Port = smap.New()
	host.Map.Tcp = smap.New()
	//host.Map.App = smap.New()

	return host
}

func (h *Host) Up() *Host {
	h.status.alive = true
	return h
}

func (h *Host) Down() *Host {
	h.status.alive = false
	return h
}

func (h *Host) SetAlivePort(port, status int) {
	h.Map.Port.Set(port, status)
	if status == Open {
		h.status.alive = true
		h.status.open = true
	}
}

func (h *Host) FinishPortScan() {
	h.status.portScan = true
	length := 0
	h.Map.Port.Range(func(key, value interface{}) bool {
		status := value.(int)
		if status != Close {
			length++
		}
		return true
	})
	h.Length.Tcp = length
}

func (h *Host) PortScanIsFinish() bool {
	return h.status.portScan
}

func (h *Host) IsAlive() bool {
	return h.status.alive
}

func (h *Host) IsOpenPort() bool {
	return h.status.open
}

func (h *Host) DisplayUnknownPorts() string {
	var portSlice []int
	var dispSlice []string

	h.Map.Tcp.Range(func(key, value interface{}) bool {
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

func (h *Host) CountUnknownPorts() int {
	i := 0
	h.Map.Tcp.Range(func(key, value interface{}) bool {
		tcpBanner := value.(*gonmap.TcpBanner)
		if tcpBanner.Status() == gonmap.Open || tcpBanner.Status() == gonmap.Unknown {
			i++
		}
		return true
	})
	return i
}

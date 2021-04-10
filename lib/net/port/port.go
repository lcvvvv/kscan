package port

import (
	"kscan/lib/misc"
)

type (
	//带协议端口
	ProtocolPorts struct {
		TCP *Ports
		UDP *Ports
	}
	//端口列表
	Ports struct {
		value  []int
		length int
	}
)

func New() *Ports {
	p := &Ports{
		value:  []int{},
		length: 0,
	}
	return p
}

func NewProtocolPorts() *ProtocolPorts {
	p := &ProtocolPorts{
		TCP: New(),
		UDP: New(),
	}
	return p
}

func (this *Ports) IsExist(i int) bool {
	if misc.IsInIntArr(this.value, i) {
		return true
	} else {
		return false
	}
}

func (this *Ports) Push(i int) bool {
	if i > 65535 || i < 0 {
		return false
	}
	if this.IsExist(i) {
		return false
	}
	this.value = append(this.value, i)
	this.length += 1
	return true
}

func (this *Ports) Pushs(iArr []int) int {
	var res int
	for _, i := range iArr {
		if this.Push(i) {
			res += 1
		}
	}
	return res
}

func (this *Ports) Len() int {
	return this.length
}

func (this *ProtocolPorts) IsExist(port int) bool {
	if this.TCP.IsExist(port) && this.UDP.IsExist(port) {
		return true
	} else {
		return false
	}
}

func (this *Ports) Load(expr string) bool {
	r := misc.MakeRegexpCompile("^(\\d+)(?:-(\\d+))?$")
	if !r.MatchString(expr) {
		return false
	}
	rArr := r.FindStringSubmatch(expr)
	var startPort, endPort int
	startPort = misc.Str2Int(rArr[1])
	if rArr[2] != "" {
		endPort = misc.Str2Int(rArr[2])
	} else {
		endPort = startPort
	}
	portArr := misc.Xrange(startPort, endPort)
	this.Pushs(portArr)
	return true
}

package run

import "fmt"

type PortStatus int

type Port struct {
	Addr   string
	Port   int
	Status PortStatus
}

func NewPort(addr string, port int) *Port {
	return &Port{
		Addr:   addr,
		Port:   port,
		Status: Unknown,
	}
}

const (
	Close   PortStatus = 0x00001a
	Open               = 0x00002b
	Unknown            = 0x00003c
)

func (p *Port) String() string {
	return fmt.Sprintf("%s:%d", p.Addr, p.Port)
}

func (p *Port) Open() *Port {
	p.Status = Open
	return p
}

func (p *Port) Close() *Port {
	p.Status = Close
	return p
}

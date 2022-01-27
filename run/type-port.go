package run

import "fmt"

type Port struct {
	addr   string
	port   int
	status int
}

func NewPort(addr string, port int) *Port {
	return &Port{
		addr:   addr,
		port:   port,
		status: Close,
	}
}

func (p *Port) SetStatus(i int) {
	p.status = i
}

func (p *Port) Status() int {
	return p.status
}

func (p *Port) UnParse() string {
	return fmt.Sprintf("%s:%d", p.addr, p.port)
}

func (p *Port) Port() int {
	return p.port
}

func (p *Port) Addr() string {
	return p.addr
}

func (p *Port) Open() *Port {
	p.status = Open
	return p
}

func (p *Port) Close() *Port {
	p.status = Close
	return p
}

func (p *Port) Unknown() *Port {
	p.status = Unknown
	return p
}

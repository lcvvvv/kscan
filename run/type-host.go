package run

import (
	"kscan/lib/smap"
)

type Host struct {
	icmp   int
	tcp    int
	status int
	port   *smap.SMap
	addr   string
}

const (
	Up = iota
	Down
	Close
	Open
	Unknown
)

func NewHost(addr string) *Host {
	return &Host{
		icmp:   Down,
		tcp:    Down,
		status: Down,
		port:   smap.New(),
		addr:   addr,
	}
}

func (h *Host) Status() int {
	return h.status
}

func (h *Host) IsOpenPort() bool {
	if h.tcp == Up {
		return true
	} else {
		return false
	}
}

func (h *Host) Up() *Host {
	h.status = Up
	return h
}

func (h *Host) Down() *Host {
	h.status = Down
	return h
}

func (h *Host) Length() int {
	return h.port.Length()
}

func (h *Host) SetIcmp(i int) {
	h.icmp = i
	if i == Up {
		h.status = Up
	}
}

func (h *Host) SetPort(port, status int) {
	h.port.Set(port, status)
	if status == Open {
		h.status = Up
		h.tcp = Up
	}
}

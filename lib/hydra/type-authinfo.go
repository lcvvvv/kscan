package hydra

import (
	"fmt"
	"kscan/lib/color"
)

type AuthInfo struct {
	Protocol string
	Port     int
	IPAddr   string
	Auth     Auth
	Status   bool
}

func NewAuthInfo(IPAddr string, Port int, Protocol string) *AuthInfo {
	a := &AuthInfo{
		Protocol: Protocol,
		Port:     Port,
		IPAddr:   IPAddr,
	}
	a.Auth = NewAuth()
	a.Status = false
	return a
}

func (a *AuthInfo) Display() string {
	s := a.Output()
	s = color.Red(s)
	s = color.Overturn(s)
	return s
}

func (a *AuthInfo) Output() string {
	s := fmt.Sprintf("%s://%s:%d\t200\tUsername:%s、Password:%s", a.Protocol, a.IPAddr, a.Port, a.Auth.Username, a.Auth.Password)
	return s
}

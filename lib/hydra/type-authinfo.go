package hydra

import (
	"fmt"
	"kscan/lib/misc"
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

func (a *AuthInfo) Output() string {
	return fmt.Sprintf("%s://%s:%d\t200\tUsername:%s、Password:%s、%s", a.Protocol, a.IPAddr, a.Port, a.Auth.Username, a.Auth.Password, misc.SprintStringMap(a.Auth.Other))
}

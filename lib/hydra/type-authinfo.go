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
	URL := fmt.Sprintf("%s://%s:%d", a.Protocol, a.IPAddr, a.Port)
	splitChar := func(i int) string {
		if i <= 23 {
			return "\t\t"
		}
		if i <= 16 {
			return "\t\t\t"
		}
		if i <= 8 {
			return "\t\t\t\t"
		}
		return "\t"
	}(len(URL))
	var s string
	if a.Auth.Username == "" {
		s = fmt.Sprintf("%s%s200\tPassword:%s", URL, splitChar, a.Auth.Password)
	} else {
		s = fmt.Sprintf("%s%s200\tUsername:%s、Password:%s", URL, splitChar, a.Auth.Username, a.Auth.Password)
	}
	s = color.Red(s)
	s = color.Overturn(s)
	return s
}

func (a *AuthInfo) Output() string {
	s := fmt.Sprintf("%s://%s:%d\t200\tUsername:%s、Password:%s", a.Protocol, a.IPAddr, a.Port, a.Auth.Username, a.Auth.Password)
	return s
}

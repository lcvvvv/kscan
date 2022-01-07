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
	authChar := ""
	if a.Auth.Username == "" {
		authChar = fmt.Sprintf("Password:%s", a.Auth.Password)
	} else {
		authChar = fmt.Sprintf("Username:%s、Password:%s", a.Auth.Username, a.Auth.Password)
	}

	var s string
	s = fmt.Sprintf("%-30v %-26v %v", URL, "Success", authChar)
	s = color.Red(s)
	s = color.Overturn(s)
	return s
}

func (a *AuthInfo) Output() string {
	URL := fmt.Sprintf("%s://%s:%d", a.Protocol, a.IPAddr, a.Port)
	authChar := ""
	if a.Auth.Username == "" {
		authChar = fmt.Sprintf("Password:%s", a.Auth.Password)
	} else {
		authChar = fmt.Sprintf("Username:%s、Password:%s", a.Auth.Username, a.Auth.Password)
	}
	var s string
	s = fmt.Sprintf("%-30v %-26v %v", URL, "Success", authChar)
	return s
}

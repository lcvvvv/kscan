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
	outMap := a.Auth.Other
	if a.Auth.Username != "" {
		outMap["Username"] = a.Auth.Username
	}
	if a.Auth.Password != "" {
		outMap["Password"] = a.Auth.Password
	}
	for key, value := range outMap {
		authChar += fmt.Sprintf("%s:%s、", key, value)
	}
	authChar = authChar[:len(authChar)-3]
	var s string
	s = fmt.Sprintf("%-30v %-26v %v", URL, "Success", authChar)
	s = color.Red(s)
	s = color.Overturn(s)
	return s
}

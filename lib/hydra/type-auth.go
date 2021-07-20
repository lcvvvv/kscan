package hydra

import "strings"

type Auth struct {
	Username string
	Password string
	Other    map[string]string
}

func NewAuth() Auth {
	a := Auth{}
	a.Other = make(map[string]string)
	return a
}

func NewSpecialAuth(username, password string, kwargs ...map[string]string) Auth {
	a := Auth{
		Username: username,
		Password: password,
	}
	a.Other = make(map[string]string)
	if len(kwargs) == 1 {
		a.Other = kwargs[0]
	}
	return a
}

func (a *Auth) MakePassword() {
	if strings.Contains(a.Password, "%user%") {
		a.Password = strings.ReplaceAll(a.Password, "%user%", a.Username)
	}
}

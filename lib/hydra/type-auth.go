package hydra

import "strings"

type Auth struct {
	Username string
	Password string
}

func NewAuth() Auth {
	a := Auth{}
	return a
}

func NewSpecialAuth(username, password string) Auth {
	a := Auth{
		Username: username,
		Password: password,
	}
	return a
}

func (a *Auth) MakePassword() {
	if strings.Contains(a.Password, "%user%") {
		a.Password = strings.ReplaceAll(a.Password, "%user%", a.Username)
	}
}

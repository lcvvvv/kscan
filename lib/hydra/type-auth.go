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

func NewAuthFromPasswords(passwords []string) []Auth {
	var auths []Auth
	for _, password := range passwords {
		auths = append(auths, Auth{
			Username: "",
			Password: password,
		})
	}
	return auths
}

func NewAuthFromUsernameAndPassword(usernames, passwords []string) []Auth {
	var auths []Auth
	for _, password := range passwords {
		for _, username := range usernames {
			auths = append(auths, Auth{
				Username: username,
				Password: password,
			})
		}
	}
	return auths
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

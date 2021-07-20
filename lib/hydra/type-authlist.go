package hydra

type AuthList struct {
	Username []string
	Password []string
	Other    map[string][]string
	Special  []Auth
}

func NewAuthList() *AuthList {
	a := &AuthList{}
	a.Other = make(map[string][]string)
	a.Special = []Auth{}
	return a
}

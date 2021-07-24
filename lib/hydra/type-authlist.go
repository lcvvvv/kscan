package hydra

import "kscan/lib/misc"

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

func (a *AuthList) IsEmpty() bool {
	if len(a.Username) > 0 && len(a.Password) > 0 {
		return false
	}
	return true
}

func (a *AuthList) Merge(list *AuthList) {
	a.Username = append(a.Username, list.Username...)
	a.Password = append(a.Password, list.Password...)
	a.Special = append(a.Special, list.Special...)
	a.Username = misc.RemoveDuplicateElement(a.Username)
	a.Password = misc.RemoveDuplicateElement(a.Password)
}

func (a *AuthList) Replace(list *AuthList) {
	a.Username = list.Username
	a.Password = list.Password
	a.Special = list.Special
	a.Username = misc.RemoveDuplicateElement(a.Username)
	a.Password = misc.RemoveDuplicateElement(a.Password)
}

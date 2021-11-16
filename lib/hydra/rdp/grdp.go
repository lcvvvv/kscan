package rdp

import (
	"fmt"
	"kscan/lib/grdp"
)

func Check(ip, domain, login, password string, port int) (bool, error) {
	target := fmt.Sprintf("%s:%d", ip, port)
	err := grdp.Login(target, domain, login, password)
	//slog.Info(target, domain, login, password)
	//slog.Info(err)

	if err != nil {
		return false, err
	}
	return true, err
}

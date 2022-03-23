package rdp

import (
	"fmt"
	"kscan/lib/grdp"
)

func Check(ip, domain, login, password string, port int, protocol string) (bool, error) {
	var err error
	target := fmt.Sprintf("%s:%d", ip, port)
	if protocol == grdp.PROTOCOL_SSL {
		err = grdp.LoginForSSL(target, domain, login, password)
	} else {
		err = grdp.LoginForRDP(target, domain, login, password)
	}
	//err = grdp.Login(target, domain, login, password)
	//slog.Info(target, domain, login, password)
	//slog.Info(err)
	if err != nil {
		return false, err
	}
	return true, err
}

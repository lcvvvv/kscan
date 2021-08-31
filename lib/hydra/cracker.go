package hydra

import (
	"kscan/lib/hydra/rdp"
	"kscan/lib/hydra/ssh"
	"kscan/lib/slog"
)

func rdpCracker(i interface{}) interface{} {
	info := i.(AuthInfo)
	info.Auth.MakePassword()
	domain := "workgroup"
	if ok, err := rdp.Check(info.IPAddr, domain, info.Auth.Username, info.Auth.Password, info.Port); ok {
		if err != nil {
			slog.Debugf("rdp://%s:%s@%s:%d:%s", info.Auth.Username, info.Auth.Password, info.IPAddr, info.Port, err)
			return nil
		}
		info.Status = true
		return info
	}
	return nil
}

func sshCracker(i interface{}) interface{} {
	info := i.(AuthInfo)
	info.Auth.MakePassword()
	if ok, err := ssh.Check(info.IPAddr, info.Auth.Username, info.Auth.Password, info.Port); ok {
		if err != nil {
			slog.Debugf("ssh://%s:%s@%s:%d:%s", info.Auth.Username, info.Auth.Password, info.IPAddr, info.Port, err)
			return nil
		}
		info.Status = true
		return info
	}
	return nil
}

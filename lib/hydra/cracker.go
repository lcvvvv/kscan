package hydra

import (
	"kscan/lib/hydra/ftp"
	"kscan/lib/hydra/mongodb"
	"kscan/lib/hydra/mssql"
	"kscan/lib/hydra/mysql"
	"kscan/lib/hydra/oracle"
	"kscan/lib/hydra/postgresql"
	"kscan/lib/hydra/rdp"
	"kscan/lib/hydra/redis"
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

func mysqlCracker(i interface{}) interface{} {
	info := i.(AuthInfo)
	info.Auth.MakePassword()
	if ok, err := mysql.Check(info.IPAddr, info.Auth.Username, info.Auth.Password, info.Port); ok {
		if err != nil {
			slog.Debugf("mysql://%s:%s@%s:%d:%s", info.Auth.Username, info.Auth.Password, info.IPAddr, info.Port, err)
			return nil
		}
		info.Status = true
		return info
	}
	return nil
}

func mssqlCracker(i interface{}) interface{} {
	info := i.(AuthInfo)
	info.Auth.MakePassword()
	if ok, err := mssql.Check(info.IPAddr, info.Auth.Username, info.Auth.Password, info.Port); ok {
		if err != nil {
			slog.Debugf("mssql://%s:%s@%s:%d:%s", info.Auth.Username, info.Auth.Password, info.IPAddr, info.Port, err)
			return nil
		}
		info.Status = true
		return info
	}
	return nil
}

func redisCracker(i interface{}) interface{} {
	info := i.(AuthInfo)
	info.Auth.MakePassword()
	if ok, err := redis.Check(info.IPAddr, info.Auth.Password, info.Port); ok {
		if err != nil {
			slog.Debugf("redis://%s:%s/auth:%s,%s", info.IPAddr, info.Port, info.Auth.Password, err)
			return nil
		}
		info.Status = true
		return info
	}
	return nil
}

func ftpCracker(i interface{}) interface{} {
	info := i.(AuthInfo)
	info.Auth.MakePassword()
	if ok, err := ftp.Check(info.IPAddr, info.Auth.Username, info.Auth.Password, info.Port); ok {
		if err != nil {
			slog.Debugf("ftp://%s:%s@%s:%d:%s", info.Auth.Username, info.Auth.Password, info.IPAddr, info.Port, err)
			return nil
		}
		info.Status = true
		return info
	}
	return nil
}

func postgresqlCracker(i interface{}) interface{} {
	info := i.(AuthInfo)
	info.Auth.MakePassword()
	if ok, err := postgresql.Check(info.IPAddr, info.Auth.Username, info.Auth.Password, info.Port); ok {
		if err != nil {
			slog.Debugf("postgres://%s:%s@%s:%d:%s", info.Auth.Username, info.Auth.Password, info.IPAddr, info.Port, err)
			return nil
		}
		info.Status = true
		return info
	}
	return nil
}

func oracleCracker(i interface{}) interface{} {
	info := i.(AuthInfo)
	info.Auth.MakePassword()
	if ok, err := oracle.Check(info.IPAddr, info.Auth.Username, info.Auth.Password, info.Port); ok {
		if err != nil {
			slog.Debugf("oracle://%s:%s@%s:%d:%s", info.Auth.Username, info.Auth.Password, info.IPAddr, info.Port, err)
			return nil
		}
		info.Status = true
		return info
	}
	return nil
}

func mongodbCracker(i interface{}) interface{} {
	info := i.(AuthInfo)
	info.Auth.MakePassword()
	if ok, err := mongodb.Check(info.IPAddr, info.Auth.Username, info.Auth.Password, info.Port); ok {
		if err != nil {
			slog.Debugf("mongodb://%s:%s@%s:%d:%s", info.Auth.Username, info.Auth.Password, info.IPAddr, info.Port, err)
			return nil
		}
		info.Status = true
		return info
	}
	return nil
}

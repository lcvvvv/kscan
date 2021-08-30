package hydra

import (
	"kscan/app"
	"kscan/lib/hydra/rdp"
	"kscan/lib/misc"
	"kscan/lib/pool"
	"kscan/lib/slog"
)

type Cracker struct {
	Pool     *pool.Pool
	authList *AuthList
	authInfo *AuthInfo
	Out      chan AuthInfo
}

func NewCracker(info *AuthInfo, threads int) *Cracker {
	c := &Cracker{}
	c.Pool = pool.NewPool(threads)
	c.authInfo = info
	if misc.IsInStrArr(app.Setting.HydraMod, c.authInfo.Protocol) == false {
		c.authInfo.Protocol = app.Setting.HydraMap[c.authInfo.Port]
	}
	c.authList = func() *AuthList {
		list := DefaultAuthMap[c.authInfo.Protocol]
		if app.Setting.HydraUpdate {
			list.Merge(CustomAuthMap)
			return list
		}
		if CustomAuthMap.IsEmpty() == false {
			list.Replace(CustomAuthMap)
			return list
		}
		return list
	}()
	c.Out = make(chan AuthInfo)
	return c
}

func (c *Cracker) Run() {
	//开启输出监测
	go c.OutWatchDog()

	switch c.authInfo.Protocol {
	case "rdp":
		c.Pool.Function = rdpCracker
		//go 任务下发器
		go func() {
			for _, password := range c.authList.Password {
				for _, username := range c.authList.Username {
					if c.Pool.Done {
						c.Pool.InDone()
						return
					}
					a := NewAuth()
					a.Password = password
					a.Username = username
					c.authInfo.Auth = a
					c.Pool.In <- *c.authInfo
				}
			}
			for _, a := range c.authList.Special {
				if c.Pool.Done {
					c.Pool.InDone()
					return
				}
				c.authInfo.Auth = a
				c.Pool.In <- *c.authInfo
			}
			//关闭信道
			c.Pool.InDone()
		}()
		//开始暴力破解
		c.Pool.Run()

	case "mysql":
	case "mssql":
	case "oracle":
	case "ldap":
	case "ssh":
	case "telnet":
	case "db2":
	case "mongodb":
	case "redis":
	case "smb":
	}
}

func (c *Cracker) OutWatchDog() {
	count := 0
	var info interface{}
	for out := range c.Pool.Out {
		if out == nil {
			continue
		}
		c.Pool.Stop()
		count += 1
		info = out
	}
	if count == 1 {
		c.Out <- info.(AuthInfo)
	}
	if count > 1 {
		slog.Debugf("rdp://%s:%d,协议不支持", info.(AuthInfo).IPAddr, info.(AuthInfo).Port)
	}
	close(c.Out)
}

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

func Ok(protocol string, port int) bool {
	if misc.IsInStrArr(app.Setting.HydraProtocolArr, protocol) {
		return true
	}
	if misc.IsInIntArr(app.Setting.HydraPortArr, port) {
		return true
	}
	return false
}

var DefaultAuthMap map[string]*AuthList
var CustomAuthMap *AuthList

func InitDefaultAuthMap() {
	m := make(map[string]*AuthList)
	m = map[string]*AuthList{
		"rdp":     NewAuthList(),
		"mysql":   NewAuthList(),
		"mssql":   NewAuthList(),
		"oracle":  NewAuthList(),
		"ldap":    NewAuthList(),
		"ssh":     NewAuthList(),
		"telnet":  NewAuthList(),
		"db2":     NewAuthList(),
		"mongodb": NewAuthList(),
		"redis":   NewAuthList(),
		"smb":     NewAuthList(),
	}
	m["rdp"] = DefaultRdpList()
	DefaultAuthMap = m
}

func InitCustomAuthMap() {
	CustomAuthMap = NewAuthList()
	CustomAuthMap.Password = app.Setting.HydraPass
	CustomAuthMap.Username = app.Setting.HydraUser
}

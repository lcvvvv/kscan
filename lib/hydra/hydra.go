package hydra

import (
	"kscan/lib/misc"
	"kscan/lib/pool"
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
	c.fixProtocol()
	c.authList = DefaultAuthMap[c.authInfo.Protocol]
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
			if _, ok := c.authList.Other["domain"]; ok == false {
				c.authList.Other["domain"] = []string{"workgroup"}
			}
			for _, password := range c.authList.Password {
				for _, username := range c.authList.Username {
					for _, domain := range c.authList.Other["domain"] {
						if c.Pool.Done {
							c.Pool.InDone()
							return
						}
						a := NewAuth()
						a.Password = password
						a.Username = username
						a.Other["domain"] = domain
						c.authInfo.Auth = a
						c.Pool.In <- *c.authInfo
					}
				}
			}
			for _, a := range c.authList.Special {
				if _, ok := a.Other["domain"]; ok == false {
					a.Other["domain"] = "workgroup"
				}
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
	for out := range c.Pool.Out {
		if out == nil {
			continue
		}
		c.Pool.Stop()
		c.Out <- out.(AuthInfo)
	}
	close(c.Out)
}

func (c *Cracker) fixProtocol() {
	protocolMap := map[int]string{
		3389:  "rdp",
		3306:  "mysql",
		1433:  "mssql",
		1521:  "oracle",
		389:   "ldap",
		22:    "ssh",
		23:    "telnet",
		21:    "ftp",
		50000: "db2",
		27017: "mongodb",
		6379:  "redis",
		445:   "smb",
	}
	c.authInfo.Protocol = protocolMap[c.authInfo.Port]

}

//func rdpCracker(i interface{}) interface{} {
//	info := i.(AuthInfo)
//	info.Auth.MakePassword()
//	domain := "workgroup"
//	if _, ok := info.Auth.Other["domain"]; ok {
//		domain = info.Auth.Other["domain"]
//	}
//	if ok, _ := rdp.Check(info.IPAddr, domain, info.Auth.Username, info.Auth.Password, info.Port); ok {
//		info.Status = true
//		return info
//	}
//	return nil
//}

func rdpCracker(i interface{}) interface{} {
	info := i.(AuthInfo)
	info.Auth.MakePassword()
	//domain := "workgroup"
	//if _, ok := info.Auth.Other["domain"]; ok {
	//	domain = info.Auth.Other["domain"]
	//}
	//if ok, _ := rdp.Check(info.IPAddr, domain, info.Auth.Username, info.Auth.Password, info.Port); ok {
	//	info.Status = true
	//	return info
	//}
	return nil
}

func Ok(protocol string, port int) bool {
	protocolArr := []string{"rdp"}
	if misc.IsInStrArr(protocolArr, protocol) {
		return true
	}
	portArr := []int{3389}
	if misc.IsInIntArr(portArr, port) {
		return true
	}
	return false
}

var DefaultAuthMap map[string]*AuthList

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

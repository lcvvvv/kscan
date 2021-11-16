package hydra

import (
	"kscan/app"
	"kscan/lib/misc"
	"kscan/lib/pool"
	"kscan/lib/slog"
	"time"
)

type Cracker struct {
	Pool     *pool.Pool
	authList *AuthList
	authInfo *AuthInfo
	Out      chan AuthInfo
}

var (
	DefaultAuthMap map[string]*AuthList
	CustomAuthMap  *AuthList
)

func NewCracker(info *AuthInfo, threads int) *Cracker {
	c := &Cracker{}
	c.Pool = pool.NewPool(threads)
	c.authInfo = info
	if misc.IsInStrArr(app.Setting.HydraProtocolArr, c.authInfo.Protocol) == false {
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
	c.Pool.Interval = time.Second * 1
	return c
}

func (c *Cracker) Run() {
	//开启输出监测
	go c.OutWatchDog()

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

	switch c.authInfo.Protocol {
	case "rdp":
		c.Pool.Function = rdpCracker
	case "mysql":
		c.Pool.Function = mysqlCracker
	case "mssql":
		c.Pool.Function = mssqlCracker
	case "oracle":
		//c.Pool.Function = oracleCracker
	case "postgresql":
		c.Pool.Function = postgresqlCracker
	case "ldap":
	case "ssh":
		c.Pool.Function = sshCracker
	case "telnet":
		c.Pool.Function = telnetCracker
	case "ftp":
		c.Pool.Function = ftpCracker
	case "mongodb":
		c.Pool.Function = mongodbCracker
	case "redis":
		c.Pool.Function = redisCracker
	case "smb":
		c.Pool.Function = smbCracker
	}
	//开始暴力破解
	c.Pool.Run()
}

func InitDefaultAuthMap() {
	m := make(map[string]*AuthList)
	m = map[string]*AuthList{
		"rdp":        NewAuthList(),
		"ssh":        NewAuthList(),
		"mysql":      NewAuthList(),
		"mssql":      NewAuthList(),
		"oracle":     NewAuthList(),
		"postgresql": NewAuthList(),
		"redis":      NewAuthList(),
		"telnet":     NewAuthList(),
		"mongodb":    NewAuthList(),
		"smb":        NewAuthList(),
		"ldap":       NewAuthList(),
		//"db2":        NewAuthList(),

	}
	m["rdp"] = DefaultRdpList()
	m["ssh"] = DefaultSshList()
	m["mysql"] = DefaultMysqlList()
	m["mssql"] = DefaultMssqlList()
	m["oracle"] = DefaultPostgresqlList()
	m["postgresql"] = DefaultPostgresqlList()
	m["redis"] = DefaultRedisList()
	m["ftp"] = DefaultFtpList()
	m["mongodb"] = DefaultMongodbList()
	m["smb"] = DefaultSmbList()
	m["telnet"] = DefaultTelnetList()
	DefaultAuthMap = m
}

func InitCustomAuthMap() {
	CustomAuthMap = NewAuthList()
	CustomAuthMap.Password = app.Setting.HydraPass
	CustomAuthMap.Username = app.Setting.HydraUser
}

func Ok(protocol string, port int) bool {
	if misc.IsInStrArr(app.Setting.HydraProtocolArr, protocol) {
		return true
	}
	//if misc.IsInIntArr(app.Setting.HydraPortArr, port) {
	//	return true
	//}
	return false
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
		slog.Debugf("%s://%s:%d,协议不支持", info.(AuthInfo).Protocol, info.(AuthInfo).IPAddr, info.(AuthInfo).Port)
	}
	close(c.Out)
}

package hydra

import (
	"errors"
	"github.com/lcvvvv/pool"
	"kscan/core/hydra/oracle"
	"kscan/lib/gotelnet"
	"kscan/lib/misc"
	"strings"
	"sync/atomic"
	"time"
)

type Cracker struct {
	Pool     *pool.Pool
	authList *AuthList
	authInfo *AuthInfo

	retries int

	SuccessCount int32
	SuccessAuth  Auth
}

func init() {
	InitDefaultAuthMap()
}

func InitDefaultAuthMap() {
	m := make(map[string]*AuthList)
	m["rdp"] = DefaultRdpList()
	m["ssh"] = DefaultSshList()
	m["mysql"] = DefaultMysqlList()
	m["mssql"] = DefaultMssqlList()
	m["oracle"] = DefaultOracleList()
	m["postgresql"] = DefaultPostgresqlList()
	m["redis"] = DefaultRedisList()
	m["ftp"] = DefaultFtpList()
	m["mongodb"] = DefaultMongodbList()
	m["smb"] = DefaultSmbList()
	m["telnet"] = DefaultTelnetList()
	DefaultAuthMap = m
}

type loginType string

const (
	ProtocolInvalid                 loginType = "ProtocolInvalid"
	UsernameAndPassword                       = "UsernameAndPassword"
	OnlyPassword                              = "OnlyPassword"
	UnauthorizedAccessVulnerability           = "UnauthorizedAccessVulnerability"
)

var (
	DefaultAuthMap map[string]*AuthList
	CustomAuthMap  *AuthList
	ProtocolList   = []string{
		"ssh", "rdp", "ftp", "smb", "telnet",
		"mysql", "mssql", "oracle", "postgresql", "mongodb", "redis",
		//110:   "pop3",
		//995:   "pop3",
		//25:    "smtp",
		//994:   "smtp",
		//143:   "imap",
		//993:   "imap",
		//389:   "ldap",
		//23:   "telnet",
		//50000: "db2",
	}
	LoginFailedErr = errors.New("login failed")
	ProtocolErr    = errors.New("protocol error")
)

func NewCracker(info *AuthInfo, isAuthUpdate bool, threads int) *Cracker {
	c := &Cracker{}
	c.retries = 3
	c.Pool = pool.New(threads)
	c.authInfo = info
	c.authList = func() *AuthList {
		list := DefaultAuthMap[c.authInfo.Protocol]
		if isAuthUpdate {
			list.Merge(CustomAuthMap)
			return list
		}
		if CustomAuthMap.IsEmpty() == false {
			list.Replace(CustomAuthMap)
			return list
		}
		return list
	}()
	c.Pool.Interval = time.Second * 1
	return c
}

func (c *Cracker) Retries(i int) {
	if i <= 0 {
		return
	}
	c.retries = i
}

func (c *Cracker) success(info Auth) {
	c.SuccessAuth = info
	atomic.AddInt32(&c.SuccessCount, 1)
}

func (c *Cracker) Run() (*Auth, error) {
	switch c.initJobFunc() {
	case ProtocolInvalid:
		return nil, ProtocolErr
	case UsernameAndPassword:
		go c.dispatcher(false)
	case OnlyPassword:
		go c.dispatcher(true)
	case UnauthorizedAccessVulnerability:
		return &UnauthorizedAccessVulnerabilityAuth, nil
	}
	//开始暴力破解
	c.Pool.Run()
	switch {
	case c.SuccessCount == 0:
		return nil, LoginFailedErr
	case c.SuccessCount <= 3:
		return &c.SuccessAuth, nil
	default:
		return nil, ProtocolErr
	}
}

func (c *Cracker) initJobFunc() loginType {
	ip := c.authInfo.IPAddr
	port := c.authInfo.Port
	//选择暴力破解函数
	switch c.authInfo.Protocol {
	case "rdp":
		c.Pool.Function = c.generateWorker(rdpCracker(ip, port))
	case "mysql":
		c.Pool.Function = c.generateWorker(mysqlCracker)
	case "mssql":
		c.Pool.Function = c.generateWorker(mssqlCracker)
	case "oracle":
		//若SID未知，则不进行后续暴力破解
		sid := oracle.GetSID(ip, port, oracle.ServiceName)
		if sid == "" {
			return ProtocolInvalid
		}
		c.Pool.Function = c.generateWorker(oracleCracker(sid))
	case "postgresql":
		c.Pool.Function = c.generateWorker(postgresqlCracker)
	case "ssh":
		c.Pool.Function = c.generateWorker(sshCracker)
	case "telnet":
		serverType := getTelnetServerType(ip, port)
		if serverType == gotelnet.UnauthorizedAccess {
			auth := NewAuth()
			auth.Other["Status"] = "UnauthorizedAccess"
			c.success(auth)
			return ProtocolInvalid
		}
		c.Pool.Function = c.generateWorker(telnetCracker(serverType))
		if serverType == gotelnet.OnlyPassword {
			return OnlyPassword
		}
	case "ftp":
		c.Pool.Function = c.generateWorker(ftpCracker)
	case "mongodb":
		c.Pool.Function = c.generateWorker(mongodbCracker)
	case "redis":
		c.Pool.Function = c.generateWorker(redisCracker)
		return OnlyPassword
	case "smb":
		c.Pool.Function = c.generateWorker(smbCracker)
	default:
		return ProtocolInvalid
	}
	return UsernameAndPassword
}

func (c *Cracker) generateWorker(f func(interface{}) error) func(interface{}) {
	return func(in interface{}) {
		for j := 0; j < c.retries; j++ {
			info := in.(AuthInfo)
			info.Auth.MakePassword()
			err := f(info)
			if err == nil {
				info.Status = true
				c.success(info.Auth)
				break
			}
			if strings.Contains(err.Error(), "timeout") == true {
				continue
			}
			if strings.Contains(err.Error(), "EOF") == true {
				continue
			}
			break
		}
	}
}

//分发器
func (c *Cracker) dispatcher(onlyPassword bool) {
	for _, auth := range c.authList.Dict(onlyPassword) {
		if c.SuccessCount > 0 {
			break
		}
		info := *c.authInfo
		info.Auth = auth
		c.Pool.Push(info)
	}
	//关闭信道
	c.Pool.Stop()
}

func InitCustomAuthMap(user, pass []string) {
	CustomAuthMap = NewAuthList()
	CustomAuthMap.Username = user
	CustomAuthMap.Password = pass
}

func Ok(protocol string) bool {
	if misc.IsDuplicate(ProtocolList, protocol) {
		return true
	}
	return false
}

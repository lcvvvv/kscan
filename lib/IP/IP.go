package IP

import (
	"../misc"
	"errors"
	"regexp"
	"strings"
)

type ip struct {
	Addr string
	Mask int
}

var regIsIP, _ = regexp.Compile("^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}$")
var regIsIPs, _ = regexp.Compile("^(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})/(\\d{1,2})$")

func check(v string) (ip, bool) {
	if regIsIP.MatchString(v) {
		vArr := strings.Split(v, ".")
		if len(vArr) != 4 {
			return ip{}, false
		}
		vIntArr, err := misc.StrArr2IntArr(vArr)
		if err != nil {
			return ip{}, false
		}
		for _, vInt := range vIntArr {
			if vInt < 0 || vInt > 255 {
				return ip{}, false
			}
		}
		return ip{v, 32}, true
	}
	if regIsIPs.MatchString(v) {
		addr := regIsIPs.FindStringSubmatch(v)[1]
		mask := regIsIPs.FindStringSubmatch(v)[2]
		_, isip := check(addr)
		if !isip {
			return ip{}, false
		}
		maskInt := misc.Str2Int(mask)
		if maskInt < 16 || maskInt > 32 {
			return ip{}, false
		}
		return ip{addr, maskInt}, true
	}
	return ip{}, false
}

func IPMask2IPArr(v string) ([]string, error) {
	IP, isIP := check(v)
	if !isIP {
		return nil, errors.New("IP格式不正确")
	}
	o := NewIpRangeLib()
	result, err := o.IpRangeToIpList(IP.Addr + "/" + misc.Int2Str(IP.Mask))
	if err != nil {
		return nil, err
	}
	if result[0][len(result[0])-2:] == ".0" {
		result = result[1:]
	}
	return result, nil
}

package IP

import (
	"errors"
	"kscan/lib/misc"
	"regexp"
	"strings"
)

type Ip struct {
	Addr string
	Mask int
}

var regIsIP, _ = regexp.Compile("^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}$")
var regIsIPs, _ = regexp.Compile("^(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})/(\\d{1,2})$")

func Check(v string) (Ip, bool) {
	if regIsIP.MatchString(v) {
		vArr := strings.Split(v, ".")
		if len(vArr) != 4 {
			return Ip{}, false
		}
		vIntArr, err := misc.StrArr2IntArr(vArr)
		if err != nil {
			return Ip{}, false
		}
		for _, vInt := range vIntArr {
			if vInt < 0 || vInt > 255 {
				return Ip{}, false
			}
		}
		return Ip{v, 32}, true
	}
	if regIsIPs.MatchString(v) {
		addr := regIsIPs.FindStringSubmatch(v)[1]
		mask := regIsIPs.FindStringSubmatch(v)[2]
		_, isip := Check(addr)
		if !isip {
			return Ip{}, false
		}
		maskInt := misc.Str2Int(mask)
		if maskInt < 0 || maskInt > 32 {
			return Ip{}, false
		}
		return Ip{addr, maskInt}, true
	}
	return Ip{}, false
}

func IPMask2IPArr(v string) ([]string, error) {
	IP, isIP := Check(v)
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

package IP

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

/*
	根据网络网段，获取该段所有IP
*/

var regxIsIP = regexp.MustCompile("^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}$")
var regxIsIPMask = regexp.MustCompile("^(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})/(\\d{1,2})$")
var regxIsIPRange = regexp.MustCompile("^(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})-(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})$")

func FormatCheck(ipExpr string) bool {
	if regxIsIP.MatchString(ipExpr) {
		return addrCheck(ipExpr)
	}
	if regxIsIPMask.MatchString(ipExpr) {
		ip := regxIsIPMask.FindStringSubmatch(ipExpr)[1]
		mask := regxIsIPMask.FindStringSubmatch(ipExpr)[2]
		if addrCheck(ip) == false {
			return false
		}
		if maskCheck(mask) == false {
			return false
		}
		return true
	}
	if regxIsIPRange.MatchString(ipExpr) {
		first := regxIsIPRange.FindStringSubmatch(ipExpr)[1]
		last := regxIsIPRange.FindStringSubmatch(ipExpr)[2]
		if addrCheck(first) == false {
			return false
		}
		if addrCheck(last) == false {
			return false
		}
		firstInt := addrStrToInt(first)
		lastInt := addrStrToInt(last)
		if firstInt > lastInt {
			return false
		}
		return true
	}
	return false
}

func ExprToList(ipExpr string) []string {
	var r []string
	if regxIsIP.MatchString(ipExpr) {
		return append(r, ipExpr)
	}
	if regxIsIPMask.MatchString(ipExpr) {
		ip := regxIsIPMask.FindStringSubmatch(ipExpr)[1]
		mask := regxIsIPMask.FindStringSubmatch(ipExpr)[2]
		maskInt, _ := strconv.Atoi(mask)
		ipInt := addrStrToInt(ip)
		maskhead := uint32(0xFFFFFFFF)
		for i := 1; i <= 32-maskInt; i++ {
			maskhead = maskhead << 1
		}
		masktail := uint32(0xFFFFFFFF)
		for i := 1; i <= maskInt; i++ {
			masktail = masktail >> 1
		}
		ipStart := uint32(ipInt) & maskhead
		ipEnd := uint32(ipInt) | masktail
		return RangeToList(ipStart, ipEnd)
	}
	if regxIsIPRange.MatchString(ipExpr) {
		start := regxIsIPRange.FindStringSubmatch(ipExpr)[1]
		end := regxIsIPRange.FindStringSubmatch(ipExpr)[2]
		startInt := addrStrToInt(start)
		endInt := addrStrToInt(end)
		return RangeToList(uint32(startInt), uint32(endInt))
	}
	return r
}

func RangeToList(start uint32, end uint32) (result []string) {
	for i := start; i <= end; i++ {
		result = append(result, addrIntToStr(int(i)))
	}
	return result
}

func addrCheck(ip string) bool {
	sArr := strings.Split(ip, ".")
	if len(sArr) != 4 {
		return false
	}
	for _, s := range sArr {
		i, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		if i > 255 || i < 0 {
			return false
		}
	}
	return true
}

func addrStrToInt(ipStr string) int {
	ipArr := strings.Split(ipStr, ".")
	var ipInt int
	var pos uint = 24
	for _, ipSeg := range ipArr {
		tempInt, _ := strconv.Atoi(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return ipInt
}

func addrIntToStr(ipInt int) string {
	ipArr := make([]string, 4)
	length := len(ipArr)
	buffer := bytes.NewBufferString("")
	for i := 0; i < length; i++ {
		tempInt := ipInt & 0xFF
		ipArr[length-i-1] = strconv.Itoa(tempInt)
		ipInt = ipInt >> 8
	}
	for i := 0; i < length; i++ {
		buffer.WriteString(ipArr[i])
		if i < length-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}

func maskCheck(mask string) bool {
	maskInt, err := strconv.Atoi(mask)
	if err != nil {
		return false
	}
	if maskInt > 32 || maskInt < 0 {
		return false
	}
	return true
}

package uri

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// IsPort checks if a string represents a valid port
func IsPort(str string) bool {
	if i, err := strconv.Atoi(str); err == nil && i > 0 && i < 65536 {
		return true
	}
	return false
}

func ParsePort(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func ParseNetlocPort(str string) (string, int) {
	r := strings.Split(str, ":")
	if len(r) != 2 {
		return "", 0
	}
	return r[0], ParsePort(r[1])
}

// IsIPv4 IsIP checks if a string is either IP version 4 Alias for `net.ParseIP`
func IsIPv4(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] == '.' {
			return net.ParseIP(str) != nil
		}
	}
	return false
}

// IsIPv6 IsIP checks if a string is either IP version 4 Alias for `net.ParseIP`
func IsIPv6(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] == ':' {
			return net.ParseIP(str) != nil
		}
	}
	return false
}

var (
	//domainRoot = []string{
	//	"com", "net", "org", "aero", "biz", "coop", "info", "museum", "name", "pro", "top", "xyz",
	//	"loan", "wang", "vip", "eu", "edu", "tech", "cloud", "online", "nrw", "cyou", "dev", "app",
	//	"shop",
	//	//country
	//	"ad", "ae", "af", "ag", "ai", "al", "am", "an", "ao", "aq", "ar", "as", "at",
	//	"au", "aw", "az", "ba", "bb", "bd", "be", "bf", "bg", "bh", "bi", "bj", "bm", "bn", "bo", "br",
	//	"bs", "bt", "bv", "bw", "by", "bz", "ca", "cc", "cf", "cg", "ch", "ci", "ck", "cl", "cm", "cn",
	//	"co", "cq", "cr", "cu", "cv", "cx", "cy", "cz", "de", "dj", "dk", "dm", "do", "dz", "ec", "ee",
	//	"eg", "eh", "es", "et", "ev", "fi", "fj", "fk", "fm", "fo", "fr", "ga", "gb", "gd", "ge", "gf",
	//	"gh", "gi", "gl", "gm", "gn", "gp", "gr", "gt", "gu", "gw", "gy", "hk", "hm", "hn", "hr", "ht",
	//	"hu", "id", "ie", "il", "in", "io", "iq", "ir", "is", "it", "jm", "jo", "jp", "ke", "kg", "kh",
	//	"ki", "km", "kn", "kp", "kr", "kw", "ky", "kz", "la", "lb", "lc", "li", "lk", "lr", "ls", "lt",
	//	"lu", "lv", "ly", "ma", "mc", "me", "md", "mg", "mh", "ml", "mm", "mn", "mo", "mp", "mq", "mr",
	//	"ms", "mt", "mv", "mw", "mx", "my", "mz", "na", "nc", "ne", "nf", "ng", "ni", "nl", "no", "np",
	//	"nr", "nt", "nu", "nz", "om", "pa", "pe", "pf", "pg", "ph", "pk", "pl", "pm", "pn", "pr", "pt",
	//	"pw", "py", "qa", "re", "ro", "ru", "rw", "sa", "sb", "sc", "sd", "se", "sg", "sh", "si", "sj",
	//	"sk", "sl", "sm", "sn", "so", "sr", "st", "su", "sy", "sz", "tc", "td", "tf", "tg", "th", "tj",
	//	"tk", "tl", "tm", "tn", "to", "tp", "tr", "tt", "tv", "tw", "tz", "ua", "ug", "uk", "us", "uy",
	//	"uz", "va", "vc", "ve", "vg", "vn", "vu", "wf", "ws", "ye", "yu", "za", "zm", "zr", "zw"}
	//domainRootString = strings.Join(domainRoot, "|")
	domainRootString = `[a-z]{2,5}`
	domainRegx       = regexp.MustCompile(`^([a-zA-Z0-9][-a-zA-Z0-9]{0,62}(?:\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})*\.(?:` + domainRootString + `))$`)
)

func IsDomain(str string) bool {
	if stringContainsCTLByte(str) == true {
		return false
	}
	if ok := domainRegx.MatchString(str); ok == false {
		return false
	}
	return true
}

func IsNetloc(str string) bool {
	return IsDomain(str) || IsIPv4(str)
}

// IsCIDR checks if the string is an valid CIDR notation (IPV4)
func IsCIDR(str string) bool {
	_, _, err := net.ParseCIDR(str)
	return err == nil
}

// IsIPRanger checks if the string is an valid CIDR notation (IPV4)
func IsIPRanger(str string) bool {
	ip1, ip2 := parseIPPairs(str)
	if ip1 != nil && ip2 != nil {
		return true
	}
	return false
}

// IsNetlocPort checks if a string is [Domain or IP]:Port
func IsNetlocPort(str string) bool {
	r := strings.Split(str, ":")
	if len(r) != 2 {
		return false
	}
	netloc := r[0]
	port := r[1]
	return IsNetloc(netloc) && IsPort(port)
}

// IsDomainPort checks if a string is Domain:Port
func IsDomainPort(str string) bool {
	r := strings.Split(str, ":")
	if len(r) != 2 {
		return false
	}
	domain := r[0]
	port := r[1]
	return IsDomain(domain) && IsPort(port)
}

// IsIPPort checks if a string is IP:Port
func IsIPPort(str string) bool {
	r := strings.Split(str, ":")
	if len(r) != 2 {
		return false
	}
	ip := r[0]
	port := r[1]
	return IsIPv4(ip) && IsPort(port)
}

func IsProtocol(str string) bool {
	ok, _ := regexp.MatchString("^[-a-z0-9A-Z]{1,20}$", str)
	return ok
}

// IsURL checks if a string is :
// protocol://netloc/path
// protocol://netloc:port/path
func IsURL(str string) bool {
	if stringContainsCTLByte(str) == true {
		return false
	}
	index := strings.Index(str, "://")
	if index == -1 {
		return false
	}
	protocol := str[:index]
	if IsProtocol(protocol) == false {
		return false
	}
	str = str[index+3:]
	if IsNetloc(str) {
		return true
	}
	if IsNetlocPort(str) {
		return true
	}
	if IsHostPath(str) {
		return true
	}
	return false
}

// IsHostPath checks if a string is :
// netloc/path
// netloc:port/path
func IsHostPath(str string) bool {
	index := strings.Index(str, "/")
	if index == -1 {
		return false
	}
	str = str[:index]
	if strings.Contains(str, ":") == true {
		return IsNetlocPort(str)
	} else {
		return IsNetloc(str)

	}
}

func GetNetlocWithURL(str string) string {
	str = str[strings.Index(str, "://")+3:]
	return GetNetlocWithHostPath(str)
}

func GetNetlocWithHostPath(str string) string {
	host := strings.TrimRight(str, "/")
	index := strings.Index(str, "/")
	if index != -1 {
		host = str[:index]
	}
	return GetNetlocWithHost(host)

}

func GetNetlocWithHost(str string) string {
	return strings.Split(str, ":")[0]
}

func SplitWithNetlocPort(str string) (netloc string, port int) {
	foo := strings.Split(str, ":")
	port, _ = strconv.Atoi(foo[1])
	return foo[0], port
}

func RangerToIP(ranger string) (IPs []net.IP) {
	first, last := parseIPPairs(ranger)
	return pairsToIP(first, last)
}

func CIDRToIP(cidr string) (IPs []net.IP) {
	_, network, _ := net.ParseCIDR(cidr)
	first := FirstIP(network)
	last := LastIP(network)
	return pairsToIP(first, last)
}

func LastIP(network *net.IPNet) net.IP {
	firstIP := FirstIP(network)
	mask, _ := network.Mask.Size()
	size := math.Pow(2, float64(32-mask))
	lastIP := toIP(toInt(firstIP) + uint32(size) - 1)
	return net.ParseIP(lastIP)
}

func FirstIP(network *net.IPNet) net.IP {
	return network.IP
}

func Contains(network *net.IPNet, ip net.IP) bool {
	i := toInt(ip)
	return toInt(FirstIP(network)) < i && toInt(LastIP(network)) > i
}

func SameSegment(ips ...string) bool {
	if len(ips) == 0 {
		return true
	}
	first := ips[0]
	_, network, _ := net.ParseCIDR(first + "/24")
	for _, ip := range ips[1:] {
		if Contains(network, net.ParseIP(ip)) == false {
			return false
		}
	}
	return true
}

// IsIPRanger parse the string is an ip pairs
// 192.168.0.1-192.168.2.255
// 192.168.0.1-255
// 192.168.0.1-2.255
func parseIPPairs(str string) (ip1 net.IP, ip2 net.IP) {
	if strings.Count(str, "-") != 1 {
		return nil, nil
	}
	r := strings.Split(str, "-")
	s1 := r[0]
	s2 := r[1]
	if ip1 = net.ParseIP(s1); ip1 == nil {
		return nil, nil
	}
	i := strings.Count(s2, ".")
	if i > 3 {
		return ip1, nil
	}
	rs1 := strings.Split(s1, ".")
	rs2 := strings.Join(append(rs1[:3-i], s2), ".")
	ip2 = net.ParseIP(rs2)
	if ip2 == nil {
		return ip1, nil
	}
	if toInt(ip1) >= toInt(ip2) {
		return nil, nil
	}
	return ip1, ip2
}

func pairsToIP(ip1, ip2 net.IP) (IPs []net.IP) {
	start := toInt(ip1)
	end := toInt(ip2)
	for i := start; i <= end; i++ {
		IPs = append(IPs, net.ParseIP(toIP(i)))
	}
	return IPs
}

// IPToInteger converts an IP address to its integer representation.
// It supports both IPv4
func toInt(ip net.IP) uint32 {
	var buf = []byte(ip)
	if len(buf) > 12 {
		buf = buf[12:]
	}
	buffer := bytes.NewBuffer(buf)
	var i uint32
	_ = binary.Read(buffer, binary.BigEndian, &i)
	return i
}

func toIP(i uint32) string {
	buf := bytes.NewBuffer([]byte{})
	_ = binary.Write(buf, binary.BigEndian, i)
	b := buf.Bytes()
	return fmt.Sprintf("%v.%v.%v.%v", b[0], b[1], b[2], b[3])
}

func GetGatewayList(ip string, t string) []string {
	var gatewayArr []string
	strArr := strings.Split(ip, ".")
	if t == "b" {
		for i := 0; i < 255; i++ {
			gatewayArr = append(gatewayArr, fmt.Sprintf("%s.%s.%d.1", strArr[0], strArr[1], i))
			gatewayArr = append(gatewayArr, fmt.Sprintf("%s.%s.%d.255", strArr[0], strArr[1], i))
		}
	}
	if t == "a" {
		for i := 0; i < 255; i++ {
			for j := 0; j < 255; j++ {
				gatewayArr = append(gatewayArr, fmt.Sprintf("%s.%d.%d.1", strArr[0], i, j))
				gatewayArr = append(gatewayArr, fmt.Sprintf("%s.%d.%d.255", strArr[0], i, j))
			}
		}
	}
	if t == "s" {
		for i := 0; i < 255; i++ {
			gatewayArr = append(gatewayArr, fmt.Sprintf("%d.%d.%d.1", i, i, i))
			gatewayArr = append(gatewayArr, fmt.Sprintf("%d.%d.%d.255", i, i, i))
		}
	}
	return gatewayArr
}

// stringContainsCTLByte reports whether s contains any ASCII control character.
func stringContainsCTLByte(s string) bool {
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b < ' ' || b == 0x7f {
			return true
		}
	}
	return false
}

func URLParse(URLRaw string) *url.URL {
	URL, _ := url.Parse(URLRaw)
	return URL
}

func GetURLPort(URL *url.URL) string {
	if port := URL.Port(); port != "" {
		return port
	}
	switch URL.Scheme {
	case "http":
		return "80"
	case "https":
		return "443"
	default:
		return ""
	}

}

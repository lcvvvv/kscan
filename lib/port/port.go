package port

import (
	"net"
)

//func IsOpen(s string) bool {
//	conn, err := net.DialTimeout("tcp", s, time.Second*2)
//	if err != nil {
//		return false
//	} else {
//		conn.Close()
//		return true
//	}
//}

func GetIP(s string) string {
	IP, err := net.ResolveIPAddr("ip4", s)
	if err != nil {
		return s
	}
	return IP.String()
}

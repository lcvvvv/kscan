package tcpping

import (
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"
)

var CommonPorts = []int{22, 23, 80, 139, 512, 443, 445, 3389}

func Ping(addr string, port int, timeout time.Duration) error {
	return connect("tcp", addr, port, timeout)
}

func PingPorts(ip string, timeout time.Duration) (err error) {
	for _, port := range CommonPorts {
		if err = Ping(ip, port, timeout); err == nil {
			return nil
		}
	}
	return err
}

func connect(protocol string, addr string, port int, duration time.Duration) error {
	host := fmt.Sprintf("%s:%d", addr, port)
	conn, err := net.DialTimeout(protocol, host, duration)
	if err != nil {
		return err
	}
	defer conn.Close()
	if (runtime.GOOS == "windows" && (port == 110 || port == 25)) == false {
		return nil
	}
	err = conn.SetDeadline(time.Now())
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		if strings.Contains(err.Error(), "forcibly closed by the remote host") {
			return err
		}
		if strings.Contains(err.Error(), "timeout") {
			return err
		}
	}
	return nil
}

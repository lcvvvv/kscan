package stcp

import (
	"io"
	"net"
	"strings"
	"time"
)

func Send(protocol string, netloc string, data string, duration time.Duration) (string, error) {
	protocol = strings.ToLower(protocol)
	conn, err := net.DialTimeout(protocol, netloc, duration)
	if err != nil {
		return "", err
	}
	buf := make([]byte, 1024)
	_, err = io.WriteString(conn, data)
	if err != nil {
		return "", err
	}
	length, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	_ = conn.Close()
	return string(buf[:length]), nil
}

package stcp

import (
	"io"
	"net"
	"time"
)

func Send(netloc string, data string, duration time.Duration) (string, error) {
	conn, err := net.DialTimeout("tcp", netloc, duration)
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

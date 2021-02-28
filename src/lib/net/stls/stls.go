package stls

import (
	"crypto/tls"
	"io"
	"net"
	"time"
)

func Send(netloc string, data string, duration time.Duration) (string, error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
	}
	dial := &net.Dialer{
		Timeout:  duration,
		Deadline: time.Now().Add(duration * 2),
	}
	conn, err := tls.DialWithDialer(dial, "tcp", netloc, config)
	if err != nil {
		return "", err
	}
	_, err = io.WriteString(conn, data)
	if err != nil {
		return "", err
	}
	buf := make([]byte, 1024)
	length, err := conn.Read(buf)
	if err != nil {
		return "", nil
	}
	_ = conn.Close()
	return string(buf[:length]), nil
}

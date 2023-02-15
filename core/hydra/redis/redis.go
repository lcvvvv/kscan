package redis

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

func Check(Host, Password string, Port int) error {
	netloc := fmt.Sprintf("%s:%d", Host, Port)
	conn, err := net.DialTimeout("tcp", netloc, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(fmt.Sprintf("auth %s\r\n", Password)))
	time.Sleep(time.Millisecond * 500)
	if err != nil {
		return err
	}
	reply, err := readResponse(conn)
	if err != nil {
		return err
	}
	if strings.Contains(reply, "+OK") == false {
		return errors.New("login failed")
	}
	return nil
}

func readResponse(conn net.Conn) (r string, err error) {
	buf := make([]byte, 4096)
	for {
		count, err := conn.Read(buf)
		if err != nil {
			break
		}
		r += string(buf[0:count])
		if count < 4096 {
			break
		}
	}
	return r, err
}

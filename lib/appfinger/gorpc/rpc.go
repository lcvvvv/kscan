package gorpc

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

var (
	bufferV1, _ = hex.DecodeString("05000b03100000004800000001000000b810b810000000000100000000000100c4fefc9960521b10bbcb00aa0021347a00000000045d888aeb1cc9119fe808002b10486002000000")
	bufferV2, _ = hex.DecodeString("050000031000000018000000010000000000000000000500")
	bufferV3, _ = hex.DecodeString("0900ffff0000")
)

func GetHostname(host string) ([]string, error) {
	netloc := fmt.Sprintf("%s:%v", host, 135)
	response, err := getResponse(netloc)
	if err != nil {
		return []string{}, err
	}
	hostname := resolvingResponse(response)
	if len(hostname) == 0 {
		return []string{}, errors.New("empty")
	}
	return hostname, nil
}

func getResponse(netloc string) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", netloc, time.Duration(5)*time.Second)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()
	err = conn.SetDeadline(time.Now().Add(time.Duration(5) * time.Second))
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(bufferV1)
	if err != nil {
		return nil, err
	}
	reply := make([]byte, 4096)
	_, err = conn.Read(reply)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(bufferV2)
	if err != nil {
		return nil, err
	}
	if n, err := conn.Read(reply); err != nil || n < 42 {
		return nil, err
	}
	text := reply[42:]
	flag := true
	for i := 0; i < len(text)-5; i++ {
		if bytes.Equal(text[i:i+6], bufferV3) {
			text = text[:i-4]
			flag = false
			break
		}
	}
	if flag {
		return nil, err
	}
	return text, nil
}

func resolvingResponse(text []byte) []string {
	encodedStr := hex.EncodeToString(text)
	express := strings.Replace(encodedStr, "0700", "", -1)
	var hostnames []string
	for _, value := range strings.Split(express, "000000") {
		value = strings.Replace(value, "00", "", -1)
		hostname, err := hex.DecodeString(value)
		if err != nil {
			continue
		}
		hostnames = append(hostnames, string(hostname))
	}
	return hostnames
}

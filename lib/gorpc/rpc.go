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

func GetFinger(host string) (*Finger, error) {
	netloc := fmt.Sprintf("%s:%v", host, 135)
	response, err := getResponse(netloc)
	if err != nil {
		return nil, err
	}
	finger := resolvingResponse(response)
	if finger.Len() == 0 {
		return nil, errors.New("finger is empty")
	}
	return finger, err
}

func getResponse(netloc string) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", netloc, time.Duration(5)*time.Second)
	if err != nil {
		return nil, err
	}
	defer func() { conn.Close() }()
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

func resolvingResponse(text []byte) *Finger {
	var finger = New()
	encodedStr := hex.EncodeToString(text)
	hostnames := strings.Replace(encodedStr, "0700", "", -1)
	hostname := strings.Split(hostnames, "000000")
	for i := 0; i < len(hostname); i++ {
		hostname[i] = strings.Replace(hostname[i], "00", "", -1)
		host, err := hex.DecodeString(hostname[i])
		if err != nil {
			return finger
		}
		finger.Append(string(host))
	}
	return finger
}

package touch

import (
	"kscan/core/gonmap"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Status bool
	Text   string
	Length int
}

func Touch(netloc string) *gonmap.TcpBanner {
	s := strings.Split(netloc, ":")
	host := s[0]
	port, _ := strconv.Atoi(s[1])
	tcpBanner := gonmap.GetTcpBanner(host, port, gonmap.New(), 3*time.Minute)
	return tcpBanner
}

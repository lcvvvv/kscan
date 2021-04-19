package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	//fmt.Println(httpfinger.FaviconHash)
	pinger, err := ping.NewPinger("10.158.4.1")
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	pinger.Timeout = time.Second * 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		panic(err)
	}
	stats := pinger
	fmt.Println(stats.PacketsRecv)
}

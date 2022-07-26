package main

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	//fmt.Println(httpfinger.FaviconHash)
	//pinger, err := ping.NewPinger("10.158.4.1")
	//if err != nil {
	//	slog.Println(slog.ERROR, err)
	//}
	//pinger.Count = 3
	//pinger.Timeout = time.Second * 3
	//err = pinger.Run() // Blocks until finished.
	//if err != nil {
	//	slog.Println(slog.ERROR, err)
	//}
	//stats := pinger
	//fmt.Println(stats.PacketsRecv)
}

type AppBanner struct {
	TcpFinger string
	AppFinger string
	Response  string
	Status    string
}

func (a *AppBanner) set(s string) {
	a.TcpFinger = s
	fmt.Println(a.TcpFinger)
}

func TestType(t *testing.T) {
	a := AppBanner{
		TcpFinger: "",
		AppFinger: "",
		Response:  "",
		Status:    "",
	}
	fmt.Println(a.TcpFinger)

	a.set("asdfadsfadsfasdf")

	fmt.Println(a.TcpFinger)

}

func TestDns(t *testing.T) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 10 * time.Second,
			}
			return d.DialContext(ctx, "udp", "114.114.114.114:513")
		},
	}

	ips, err := r.LookupHost(context.Background(), "asdfasdfasfasdf")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ips)
}

func TestContext1(t *testing.T) {
	go Context()
	time.Sleep(10 * time.Second)
}

func Context() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resChan := make(chan string)

	go func() {
		defer func() {
			if err := recover(); err != nil {
			}
		}()
		time.Sleep(3 * time.Second)
		resChan <- ""
	}()

	for {
		select {
		case <-ctx.Done():
			close(resChan)
			fmt.Println("超时执行结束")
			return
		case <-resChan:
			fmt.Println("正常执行完毕")
			close(resChan)
			return
		}
	}
}

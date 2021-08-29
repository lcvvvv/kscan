package main

import (
	"context"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	//fmt.Println(httpfinger.FaviconHash)
	//pinger, err := ping.NewPinger("10.158.4.1")
	//if err != nil {
	//	slog.Error(err)
	//}
	//pinger.Count = 3
	//pinger.Timeout = time.Second * 3
	//err = pinger.Run() // Blocks until finished.
	//if err != nil {
	//	slog.Error(err)
	//}
	//stats := pinger
	//fmt.Println(stats.PacketsRecv)

	utf8Str := "编码转换内容内容"
	utf8Buf := []byte(utf8Str)

	gbkBuf, _ := simplifiedchinese.GBK.NewEncoder().Bytes(utf8Buf)
	gbkStr := string(gbkBuf)
	fmt.Println(gbkBuf) //byte
	fmt.Println(gbkStr) //打印为乱码
	utf8Buf, _ = simplifiedchinese.GBK.NewDecoder().Bytes(gbkBuf)
	utf8Str = string(utf8Buf)
	fmt.Println(utf8Buf) //byte
	fmt.Println(utf8Str) //打印为乱码
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

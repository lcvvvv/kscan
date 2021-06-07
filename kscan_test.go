package main

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"testing"
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

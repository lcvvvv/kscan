package stdio

import (
	"github.com/lcvvvv/stdio/chinese"
	"os"
	"runtime"
)

var (
	Out = &Writer{}
	Err = &Writer{}
)

const (
	utf8 = 0x00000a1
	gbk  = 0x00000b2
)

var encoding = utf8

type Writer struct{}

func init() {
	if runtime.GOOS == "windows" {
		encoding = gbk
	}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	var b []byte
	if encoding == utf8 {
		b = chinese.ByteToUTF8(p)
		//fmt.Println("UTF-8")
	} else {
		b = chinese.ByteToGBK(p)
		//fmt.Println("GBK")
	}
	return os.Stdout.Write(b)
}

func SetUTF8() {
	encoding = utf8
}

func SetGBK() {
	encoding = gbk
}

func SetEncoding(s string) {
	if s == "" {
		return
	}
	if s == "gbk" {
		encoding = gbk
		return
	}
	if s == "gb2312" {
		encoding = gbk
		return
	}
	if s == "utf8" {
		encoding = utf8
		return
	}
	if s == "utf-8" {
		encoding = utf8
		return
	}
	panic("encoding value is invalid")
}

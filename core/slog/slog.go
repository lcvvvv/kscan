package slog

import (
	"fmt"
	"github.com/lcvvvv/gonmap/lib/chinese"
	"io"
	"io/ioutil"
	"kscan/lib/color"
	"log"
	"os"
	"runtime"
	"strings"
)

var splitStr = "> "

type Level int

type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}

type logger struct {
	log      *log.Logger
	modifier func(string) string
}

func (l *logger) Printf(format string, s ...interface{}) {
	expr := fmt.Sprintf(format, s...)
	l.Println(expr)
}

func (l *logger) Println(s ...interface{}) {
	expr := fmt.Sprint(s...)
	if l.modifier != nil {
		expr = l.modifier(expr)
	}
	l.log.Println(expr)
}

var info = Logger(log.New(ioutil.Discard, "", 0))
var warn = Logger(log.New(ioutil.Discard, "", 0))
var err = Logger(log.New(ioutil.Discard, "", 0))
var dbg = Logger(log.New(ioutil.Discard, "", 0))
var data = Logger(log.New(os.Stdout, "\r", 0))

func SetEncoding(v string) {
	encoding = v
}

var encoding = "utf-8"

const (
	DEBUG Level = 0x0000a1
	INFO        = 0x0000b2
	WARN        = 0x0000c3
	ERROR       = 0x0000d4
	DATA        = 0x0000f5
	NONE        = 0x0000e6
)

func Printf(level Level, format string, s ...interface{}) {
	Println(level, fmt.Sprintf(format, s...))
}

func Println(level Level, s ...interface{}) {
	logStr := fmt.Sprint(s...)
	if encoding == "gb2312" {
		logStr = chinese.ToGBK(logStr)
	} else {
		logStr = chinese.ToUTF8(logStr)
	}

	switch level {
	case DEBUG:
		if debugFilter(logStr) {
			return
		}
		dbg.Println(logStr)
	case INFO:
		info.Println(logStr)
	case WARN:
		warn.Println(logStr)
	case ERROR:
		err.Println(logStr)
	case DATA:
		data.Println(logStr)
	default:
		return
	}
}

func Debug() Logger {
	return dbg
}

func SetLogger(level Level) {
	if level <= ERROR {
		err = Logger(&logger{
			log.New(io.MultiWriter(os.Stderr), "\rError:", 0),
			nil,
		})
	}
	if level <= WARN {
		warn = Logger(&logger{
			log.New(os.Stdout, "\r[*]", log.Ldate|log.Ltime),
			color.Red,
		})
	}
	if level <= INFO {
		info = Logger(&logger{
			log.New(os.Stdout, "\r[+]", log.Ldate|log.Ltime),
			color.Green,
		})
	}
	if level <= DEBUG {
		dbg = Logger(&logger{
			log.New(os.Stdout, "\r[-]", log.Ldate|log.Ltime),
			debugModifier,
		})
	}
	if level <= NONE {
		//nothing
	}
}

func debugModifier(s string) string {
	_, file, line, _ := runtime.Caller(3)
	file = file[strings.LastIndex(file, "/")+1:]
	logStr := fmt.Sprintf("%s%s(%d) %s", splitStr, file, line, s)
	logStr = color.Yellow(logStr)
	return logStr
}

func debugFilter(s string) bool {
	//Debug 过滤器
	if strings.Contains(s, "too many open") { //发现存在线程过高错误
		Println(ERROR, "当前线程过高，请降低线程!或者请执行\"ulimit -n 50000\"命令放开操作系统限制,MAC系统可能还需要执行：\"launchctl limit maxfiles 50000 50000\"")
	}
	//if strings.Contains(s, "STEP1:CONNECT") {
	//	return true
	//}
	return false
}

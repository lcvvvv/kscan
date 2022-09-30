package slog

import (
	"fmt"
	"github.com/lcvvvv/stdio"
	"io"
	"io/ioutil"
	"kscan/lib/color"
	"log"
	"os"
	"runtime"
	"strings"
)

type Level int

const (
	DEBUG Level = 0x0000a1
	INFO        = 0x0000b2
	WARN        = 0x0000c3
	ERROR       = 0x0000d4
	DATA        = 0x0000f5
	NONE        = 0x0000e6
)

type Logger struct {
	log      *log.Logger
	modifier func(string) string
	filter   func(string) bool
}

func (l *Logger) Printf(format string, s ...interface{}) {
	expr := fmt.Sprintf(format, s...)
	l.Println(expr)
}

func (l *Logger) Println(s ...interface{}) {
	expr := fmt.Sprint(s...)
	if l.modifier != nil {
		expr = l.modifier(expr)
	}
	if l.filter != nil {
		if l.filter(expr) == true {
			return
		}
	}
	l.log.Println(expr)
}

var info = &Logger{
	log.New(stdio.Out, "\r[+]", log.Ldate|log.Ltime),
	color.Green,
	nil,
}

var warn = &Logger{
	log.New(stdio.Out, "\r[*]", log.Ldate|log.Ltime),
	color.Red,
	nil,
}

var err = &Logger{
	log.New(io.MultiWriter(stdio.Err), "\rError:", 0),
	nil,
	nil,
}

var dbg = &Logger{
	log.New(stdio.Out, "\r[-]", log.Ldate|log.Ltime),
	debugModifier,
	debugFilter,
}

func debugModifier(s string) string {
	_, file, line, _ := runtime.Caller(3)
	file = file[strings.LastIndex(file, "/")+1:]
	logStr := fmt.Sprintf("%s%s(%d) %s", "> ", file, line, s)
	logStr = color.Yellow(logStr)
	return logStr
}

func debugFilter(s string) bool {
	//Debug 过滤器
	if strings.Contains(s, "too many open") { //发现存在线程过高错误
		fmt.Println("当前线程过高，请降低线程!或者请执行\"ulimit -n 50000\"命令放开操作系统限制")
		os.Exit(0)
	}
	//if strings.Contains(s, "STEP1:CONNECT") {
	//	return true
	//}
	return false
}

var data = &Logger{
	log.New(stdio.Out, "\r", 0),
	nil,
	nil,
}

func Printf(level Level, format string, s ...interface{}) {
	Println(level, fmt.Sprintf(format, s...))
}

func Println(level Level, s ...interface{}) {
	logStr := fmt.Sprint(s...)
	switch level {
	case DEBUG:
		dbg.Println(logStr)
	case INFO:
		info.Println(logStr)
	case WARN:
		warn.Println(logStr)
	case ERROR:
		err.Println(logStr)
		os.Exit(0)
	case DATA:
		data.Println(logStr)
	default:
		return
	}
}

var empty = &Logger{log.New(ioutil.Discard, "", 0), nil, nil}

func SetLevel(level Level) {
	if level > ERROR {
		err = empty
	}
	if level > WARN {
		warn = empty
	}
	if level > INFO {
		info = empty
	}
	if level > DEBUG {
		dbg = empty
	}
	if level > NONE {
		//nothing
	}
}

func SetOutput(writer io.Writer) {
	data.modifier = func(s string) string {
		_, _ = writer.Write([]byte(color.Clear(s)))
		_, _ = writer.Write([]byte("\r\n"))
		return s
	}
}

func Debug() *Logger {
	return dbg
}

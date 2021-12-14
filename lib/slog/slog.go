package slog

import (
	"fmt"
	"io"
	"io/ioutil"
	"kscan/lib/chinese"
	"kscan/lib/color"
	"log"
	"os"
	"runtime"
	"strings"
)

var this logger
var splitStr = "> "

type (
	LEVEL int

	logger struct {
		info     *log.Logger
		warning  *log.Logger
		debug    *log.Logger
		data     *log.Logger
		error    *log.Logger
		encoding string
		//fooLine *log.Logger
	}
)

const (
	DEBUG LEVEL = iota
	INFO
	WARN
	ERROR
	NONE
)

func init() {
	this.info = log.New(ioutil.Discard, "", 0)
	this.warning = log.New(ioutil.Discard, "", 0)
	this.error = log.New(ioutil.Discard, "", 0)
	this.debug = log.New(ioutil.Discard, "", 0)
	this.data = log.New(os.Stdout, "\r", 0)
}

func SetEncoding(encoding string) {
	this.encoding = encoding
}

func SetPrintDebug(PrintDebug bool) {
	if PrintDebug {
		SetLogger(DEBUG)
	} else {
		SetLogger(INFO)
	}
}

func SetLogger(level LEVEL) {
	if level <= ERROR {
		this.error = log.New(io.MultiWriter(os.Stderr), "\rError:", 0)
	}
	if level <= WARN {
		this.warning = log.New(os.Stdout, "\r[*]", log.Ldate|log.Ltime)
	}
	if level <= INFO {
		this.info = log.New(os.Stdout, "\r[+]", log.Ldate|log.Ltime)
	}
	if level <= DEBUG {
		this.debug = log.New(os.Stdout, "\r[-]", log.Ldate|log.Ltime)
	}
	if level <= NONE {
		//nothing
	}
}

func (t *logger) Data(s string) {
	t.data.Print(s)
}

func (t *logger) Info(s string) {
	s = color.Green(s)
	t.info.Print(splitStr, s)
}

func (t *logger) Error(s string) {
	t.error.Print(s)
	os.Exit(0)
}

func (t *logger) Warning(s string) {
	s = color.Red(s)
	t.warning.Print(splitStr, s)
}

func (t *logger) Debug(s string) {
	if debugFilter(s) {
		return
	}
	_, file, line, _ := runtime.Caller(3)
	file = file[strings.LastIndex(file, "/")+1:]
	logStr := fmt.Sprintf("%s%s(%d) %s", splitStr, file, line, s)
	logStr = color.Yellow(logStr)
	t.debug.Printf(logStr)
}

func (t *logger) DoPrint(logType string, logStr string) {
	if this.encoding == "gb2312" {
		logStr = chinese.ToGBK(logStr)
	} else {
		logStr = chinese.ToUTF8(logStr)
	}
	switch logType {
	case "Debug":
		t.Debug(logStr)
	case "Info":
		t.Info(logStr)
	case "Data":
		t.Data(logStr)
	case "Warning":
		t.Warning(logStr)
	case "Error":
		t.Error(logStr)
	}

}

//func (t *logger) Error(s string) {
//	_, file, line, _ := runtime.Caller(2)
//	file = file[strings.LastIndex(file, "/")+1:]
//	t.error.Printf("%s%s(%d) %s", splitStr, file, line, s)
//	os.Exit(0)
//}

//func (t *logger) Errorf(format string, v ...interface{}) {
//	_, file, line, _ := runtime.Caller(2)
//	file = file[strings.LastIndex(file, "/")+1:]
//	format = fmt.Sprintf("%s%s(%d) %s", splitStr, file, line, format)
//	t.error.Printf(format, v...)
//	os.Exit(0)
//}
func Error(s ...interface{}) {
	this.DoPrint("Error", fmt.Sprint(s...))
}

func Info(s ...interface{}) {
	this.DoPrint("Info", fmt.Sprint(s...))
}

func Infof(format string, v ...interface{}) {
	this.DoPrint("Info", fmt.Sprintf(format, v...))
}

func Warning(s ...interface{}) {
	this.DoPrint("Warning", fmt.Sprint(s...))
}

func Warningf(format string, v ...interface{}) {
	this.DoPrint("Warning", fmt.Sprintf(format, v...))
}

func Debug(s ...interface{}) {
	this.DoPrint("Debug", fmt.Sprint(s...))
}

func Debugf(format string, v ...interface{}) {
	this.DoPrint("Debug", fmt.Sprintf(format, v...))
}

func Data(v ...interface{}) {
	this.DoPrint("Data", fmt.Sprint(v...))
}

//func Error(s string) {
//	this.Error(s)
//}

//func Errorf(format string, v ...interface{}) {
//	this.Errorf(format, v...)
//}

func debugFilter(s string) bool {
	//Debug 过滤器
	if strings.Contains(s, "too many open") { //发现存在线程过高错误
		Error("当前线程过高，请降低线程!或者请执行\"ulimit -n 50000\"命令放开操作系统限制,MAC系统可能还需要执行：\"launchctl limit maxfiles 50000 50000\"")
	}
	//if strings.Contains(s, "STEP1:CONNECT") {
	//	return true
	//}
	return false
}

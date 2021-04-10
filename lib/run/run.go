package run

import (
	"fmt"
	"github.com/lcvvvv/gonmap"
	"kscan/app"
	"kscan/lib/misc"
	"kscan/lib/queue"
	"kscan/lib/scan"
	"kscan/lib/slog"
	"net"
	"sync"
	"time"
)

var HostPortQueue = queue.New()
var HostQueue = queue.New()

var HostNum int

//var PortQueue = queue.New()

//var threadPortSync int
var threadOpenPortGroup sync.WaitGroup
var threadOpenPortGroupNum int

func InitPortQueue() {
	for _, host := range app.Config.HostTarget {
		HostQueue.Push(host)
	}
	HostNum = HostQueue.Len()
	slog.Warningf("总共扫描主机对象%d个...", HostNum)
	go initPortQueueSub()
	time.Sleep(time.Second * 1)
}

func initPortQueueSub() {
	for _, host := range app.Config.HostTarget {
		for _, Port := range app.Config.Port {
			IP := GetIP(host)
			HostPortQueue.Push(fmt.Sprintf("%s:%d", IP, Port))
		}
		for {
			if HostPortQueue.Len() < 65535 {
				break
			}
			time.Sleep(time.Second * 1)
		}
		_ = HostQueue.Pop()
	}
}

func InitHostPortQueue() int {
	for _, url := range app.Config.UrlTarget {
		HostPortQueue.Push(url)
	}
	return len(app.Config.UrlTarget)
}

func Start() {
	var thread int
	thread = app.Config.Threads
	for i := 0; i <= thread; i++ {
		threadOpenPortGroup.Add(1)
		threadOpenPortGroupNum++
		go startSub(&HostPortQueue, &threadOpenPortGroup, nil)
	}
	go WatchDogSub(&HostPortQueue)
	threadOpenPortGroup.Wait()
}

func startSub(HostPortQueue **queue.Queue, wait *sync.WaitGroup, nmap *gonmap.Nmap) {
	if nmap == nil {
		nmap = gonmap.New()
	}
	if (*HostPortQueue).Len() > 0 {
		t := misc.Interface2Str((*HostPortQueue).Pop())
		//fmt.Printf("\r[*][%d][%d/%d][协程数：%d]正在测试端口开放情况：%s", (*HostPortQueue).Len(), HostQueue.Len(), HostNum, threadOpenPortGroupNum, t)
		r := scan.GetPortBanner(t, nmap)
		if r.Status != "CLOSED" && r.Status != "KNOWN" {
			r.MakeInfo()
			slog.Data(r.Info)
		}
	}
	if (*HostPortQueue).Len() == 0 {
		time.Sleep(time.Second * 2)
	}
	if (*HostPortQueue).Len() > 0 {
		startSub(HostPortQueue, wait, nmap)
		return
	} else {
		threadOpenPortGroupNum--
		wait.Done()
		return
	}
}

func WatchDogSub(HostPortQueue **queue.Queue) {
	for {
		length := (*HostPortQueue).Len()
		if length > 0 {
			t := (*HostPortQueue).Peek()
			line := fmt.Sprintf("[%d][%d/%d][协程数：%d]正在测试端口开放情况：%s", length, HostQueue.Len(), HostNum, threadOpenPortGroupNum, t)
			slog.FooLine(line)
		}
		if length == 0 {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func GetIP(s string) string {
	IP, err := net.ResolveIPAddr("ip4", s)
	if err != nil {
		return s
	}
	return IP.String()
}

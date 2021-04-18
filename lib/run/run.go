package run

import (
	"fmt"
	"github.com/lcvvvv/gonmap"
	"kscan/app"
	"kscan/lib/misc"
	"kscan/lib/queue"
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
var threadHostPortGroup sync.WaitGroup
var threadHostPortGroupNum int

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
	for _, Port := range app.Config.Port {
		for _, host := range app.Config.HostTarget {
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
	if HostPortQueue.Len() < thread {
		thread = HostPortQueue.Len()
	}
	for i := 0; i <= thread; i++ {
		threadHostPortGroup.Add(1)
		threadHostPortGroupNum++
		go startSub(&HostPortQueue, &threadHostPortGroup, nil)
	}
	go WatchDogSub(&HostPortQueue)
	threadHostPortGroup.Wait()
}

func startSub(HostPortQueue **queue.Queue, wait *sync.WaitGroup, nmap *gonmap.Nmap) {
	if nmap == nil {
		nmap = gonmap.New()
	}
	if (*HostPortQueue).Len() > 0 {
		t := (*HostPortQueue).Pop().(string)
		r := GetPortBanner(t, nmap)
		if r.Status == "OPEN" || r.Status == "MATCHED" {
			if r.Finger.Service != "" {
				r.MakeInfo()
				slog.Data(r.Info)
				if app.Config.Output != nil {
					_, _ = app.Config.Output.WriteString(r.Info + "\n")
				}
			}
		}
	}
	if (*HostPortQueue).Len() == 0 {
		time.Sleep(time.Millisecond * 500)
	}
	if (*HostPortQueue).Len() == 0 {
		threadHostPortGroupNum--
		wait.Done()
	} else {
		startSub(HostPortQueue, wait, nmap)
	}
}

func WatchDogSub(HostPortQueue **queue.Queue) {
	for {
		if (*HostPortQueue).Len() > 0 {
			//t := (*HostPortQueue).Peek()
			var percent string
			if HostQueue.Len() > 0 {
				percent = misc.Percent(HostQueue.Len()*len(app.Config.Port)+65535, HostNum*len(app.Config.Port))
			} else {
				percent = misc.Percent((*HostPortQueue).Len(), HostNum*len(app.Config.Port))
			}
			line := fmt.Sprintf("[%s%%][%d/%d][协程数：%d]正在测试端口开放情况情况....", percent, HostQueue.Len(), HostNum, threadHostPortGroupNum)
			slog.FooLine(line)
		}
		time.Sleep(time.Second * 2)
		if (*HostPortQueue).Len() == 0 {
			line := fmt.Sprintf("所有探针已下发完毕，目前[存活协程数：%d]...", threadHostPortGroupNum)
			slog.FooLine(line)
			if threadHostPortGroupNum == 0 {
				slog.FooLine("扫描结束，现在退出程序...")
				break
			}
		}
	}
}

func GetIP(s string) string {
	IP, err := net.ResolveIPAddr("ip4", s)
	if err != nil {
		return s
	}
	return IP.String()
}

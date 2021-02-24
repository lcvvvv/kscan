package run

import (
	"app/params"
	"fmt"
	"lib/misc"
	"lib/port"
	"lib/queue"
	"lib/slog"
	"sync"
	"time"
)

var OpenPortQueue = queue.New()
var HostQueue = queue.New()

var HostNum int

//var PortQueue = queue.New()

//var threadPortSync int
var threadOpenPortGroup sync.WaitGroup
var threadOpenPortGroupNum int

func InitPortQueue() {
	for _, host := range params.SerParams.HostTarget {
		HostQueue.Push(host)
	}
	HostNum = HostQueue.Len()
	slog.Warningf("总共扫描主机对象%d个...", HostNum)
	go InitPortQueueSub()
	time.Sleep(time.Second * 1)
}

func InitPortQueueSub() {
	for _, host := range params.SerParams.HostTarget {
		for _, Port := range params.SerParams.Port {
			IP := port.GetIP(host)
			OpenPortQueue.Push(fmt.Sprintf("%s:%d", IP, Port))
		}
		for {
			if OpenPortQueue.Len() < 65535 {
				break
			}
			time.Sleep(time.Second * 1)
		}
		_ = HostQueue.Pop()
	}
}

func InitUrlQueue() {
	for _, url := range params.SerParams.UrlTarget {
		OpenPortQueue.Push(url)
	}
}

func GetBanner() {
	var thread int
	thread = params.SerParams.Threads
	for i := 0; i <= thread; i++ {
		threadOpenPortGroup.Add(1)
		threadOpenPortGroupNum++
		go GetBannerSub(&OpenPortQueue, &threadOpenPortGroup)
	}
	go WatchDogSub(&OpenPortQueue)
	threadOpenPortGroup.Wait()
}

func GetBannerSub(OpenPortQueue **queue.Queue, wait *sync.WaitGroup) {
	if (*OpenPortQueue).Len() > 0 {
		t := misc.Interface2Str((*OpenPortQueue).Pop())
		//fmt.Printf("\r[*][%d][%d/%d][协程数：%d]正在测试端口开放情况：%s", (*OpenPortQueue).Len(), HostQueue.Len(), HostNum, threadOpenPortGroupNum, t)
		port.GetBanner(t)
	}
	if (*OpenPortQueue).Len() == 0 {
		time.Sleep(time.Second * 2)
	}
	if (*OpenPortQueue).Len() > 0 {
		GetBannerSub(OpenPortQueue, wait)
		return
	} else {
		threadOpenPortGroupNum--
		wait.Done()
		return
	}
}

func WatchDogSub(OpenPortQueue **queue.Queue) {
	for {
		length := (*OpenPortQueue).Len()
		if length > 0 {
			t := (*OpenPortQueue).Peek()
			line := fmt.Sprintf("[%d][%d/%d][协程数：%d]正在测试端口开放情况：%s", length, HostQueue.Len(), HostNum, threadOpenPortGroupNum, t)
			slog.FooLine(line)
		}
		if length == 0 {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
}

package run

import (
	"../../app/params"
	"../../lib/misc"
	"../../lib/port"
	"../../lib/queue"
	"fmt"
	"sync"
)

var OpenPortQueue = queue.New()

//var PortQueue = queue.New()

//var threadPortSync int
var threadOpenPortGroup sync.WaitGroup

func InitPortQueue() {
	for _, host := range params.SerParams.HostTarget {
		for _, Port := range params.SerParams.Port {
			IP := port.GetIP(host)
			//fmt.Print(fmt.Sprintf("%s:%d", IP, Port),"\n")
			//PortQueue.Push(fmt.Sprintf("%s:%d", IP, Port))
			OpenPortQueue.Push(fmt.Sprintf("%s:%d", IP, Port))
		}
	}
}

func InitUrlQueue() {
	for _, url := range params.SerParams.UrlTarget {
		OpenPortQueue.Push(url)
	}
}

//func ScanOpenPort() {
//	var thread int
//	if PortQueue.Len() > 200 {
//		thread = params.SerParams.Threads
//	} else {
//		thread = 20
//	}
//	fmt.Print("[*]总共线程数：", thread, "\n")
//	for i := 0; i <= thread; i++ {
//		threadPortSync++
//		go ScanOpenPortSub(&PortQueue, &OpenPortQueue, &threadPortSync)
//	}
//	fmt.Printf("[*]成功建立%d个线程探测端口开放情况\n", thread)
//}
//
//func ScanOpenPortSub(PortQueue **queue.Queue, OpenPortQueue **queue.Queue, wait *int) {
//	t := misc.Interface2Str((*PortQueue).Pop())
//	//fmt.Printf("\r[*]正在测试端口是否开放：%s",t)
//	//fmt.Print(t,"\n"
//	if port.IsOpen(t) {
//		(*OpenPortQueue).Push(t)
//	}
//	//fmt.Print((*PortQueue).Empty(),"\n")
//	if (*PortQueue).Len()>0 {
//		ScanOpenPortSub(PortQueue, OpenPortQueue, wait)
//	} else {
//		*wait--
//	}
//}

func GetBanner() {
	var thread int
	lenth := OpenPortQueue.Len()
	if lenth > 200 {
		thread = params.SerParams.Threads
	} else {
		thread = 20
	}
	for i := 0; i <= thread; i++ {
		threadOpenPortGroup.Add(1)
		go GetBannerSub(&OpenPortQueue, &threadOpenPortGroup, lenth)
	}
	threadOpenPortGroup.Wait()
}

//func GetBannerSub(OpenPortQueue **queue.Queue, scanWait *int, wait *sync.WaitGroup) {
//	if (*OpenPortQueue).Len() > 0 {
//		t := misc.Interface2Str((*OpenPortQueue).Pop())
//		port.GetBanner(t)
//	}
//	if (*OpenPortQueue).Len() > 0 {
//		GetBannerSub(OpenPortQueue, scanWait, wait)
//		return
//	} else {
//		//如果开放端口队列为空，判断端口扫描程序是否结束
//		if *scanWait != 0 {
//			//未结束则等待，之后重新判断开放端口队列是否为空
//			time.Sleep(time.Millisecond * time.Duration(2000))
//			GetBannerSub(OpenPortQueue, scanWait, wait)
//			return
//		} else {
//			//如果结束了，则扫描完毕
//			wait.Done()
//			return
//		}
//	}
//}

func GetBannerSub(OpenPortQueue **queue.Queue, wait *sync.WaitGroup, lenth int) {
	if (*OpenPortQueue).Len() > 0 {
		t := misc.Interface2Str((*OpenPortQueue).Pop())
		fmt.Printf("\r[*][%d/%d]正在测试端口开放情况：%s", (*OpenPortQueue).Len(), lenth, t)
		port.GetBanner(t)
	}
	//fmt.Print((*OpenPortQueue).Len())
	if (*OpenPortQueue).Len() > 0 {
		GetBannerSub(OpenPortQueue, wait, lenth)
		return
	} else {
		wait.Done()
		return
	}
}

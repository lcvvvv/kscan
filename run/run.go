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
var RealPortNum int

var threadHostPortGroup sync.WaitGroup
var threadHostPortGroupNum int

func Start() {
	RealPortNum = app.Config.PortNum
	//STEP0:初始化进度检测器
	go WatchDogSub()
	//STEP1:初始化可访问URL地址队列
	slog.Warning("开始压入URL地址队列...")
	go pushUrlTarget()
	slog.Warningf("本次需要直接扫描的URL地址共:[%d]个。...", app.Config.UrlTargetNum)

	//STEP2:初始化主机端口扫描队列
	slog.Warningf("本次需扫描主机IP地址共:[%d]个...", app.Config.HostTargetNum)
	go pushHostTarget()
	slog.Warning("开始压入端口扫描队列...")
	time.Sleep(time.Second * 1) //预留一秒钟加载时间

	//STEP3:开始扫描所有开放端口
	go scanMain()               //启动扫描主程序
	time.Sleep(time.Second * 1) //预留一秒钟加载时间

	//STEP4:等待程序运行结束
	threadHostPortGroup.Wait()
}

func scanMain() {
	var thread = app.Config.Threads
	if HostPortQueue.Len() < thread {
		thread = HostPortQueue.Len()
	}
	slog.Warningf("开始扫描所有开放端口,总协程数为：[%d]...", thread)
	for i := 0; i <= thread; i++ {
		threadHostPortGroup.Add(1)
		threadHostPortGroupNum++
		go scanMainSub(&HostPortQueue, &threadHostPortGroup, nil)
	}
}

func scanMainSub(HostPortQueue **queue.Queue, wait *sync.WaitGroup, nmap *gonmap.Nmap) {
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
		time.Sleep(time.Second * 5)
	}
	if (*HostPortQueue).Len() == 0 {
		threadHostPortGroupNum--
		wait.Done()
	} else {
		scanMainSub(HostPortQueue, wait, nmap)
	}
}

func GetIP(s string) string {
	IP, err := net.ResolveIPAddr("ip4", s)
	if err != nil {
		return s
	}
	return IP.String()
}

func pushUrlTarget() {
	for _, url := range app.Config.UrlTarget {
		HostPortQueue.Push(url)
	}
}

func pushHostTarget() {
	for _, Port := range app.Config.Port {
		for _, host := range app.Config.HostTarget {
			IP := GetIP(host)
			HostPortQueue.Push(fmt.Sprintf("%s:%d", IP, Port))
		}
		for {
			if HostPortQueue.Len() < 4000 {
				break
			}
			time.Sleep(time.Millisecond * 300)
		}
		RealPortNum--
	}
	slog.Warningf("所有待扫描的端口已全部压入队列...")
}

func WatchDogSub() {
	HostTargetNum := app.Config.HostTargetNum
	PortNum := app.Config.PortNum
	for {
		time.Sleep(time.Second * 10)
		if HostPortQueue.Len() > 0 {
			var percent string
			if RealPortNum > 0 {
				percent = misc.Percent(RealPortNum*HostTargetNum+4000, HostTargetNum*PortNum)
			} else {
				percent = misc.Percent(HostPortQueue.Len(), app.Config.HostTargetNum*PortNum)
			}
			line := fmt.Sprintf("[%s%%][%d/%d][协程数：%d]正在测试端口开放情况情况....", percent, RealPortNum, PortNum, threadHostPortGroupNum)
			slog.Info(line)
		}
		time.Sleep(time.Second * 10)
		if HostPortQueue.Len() == 0 {
			line := fmt.Sprintf("所有探针已下发完毕，目前[存活协程数：%d]...", threadHostPortGroupNum)
			slog.Info(line)
			if threadHostPortGroupNum == 0 {
				slog.Info("扫描结束，现在退出程序...")
				break
			}
		}
	}
}

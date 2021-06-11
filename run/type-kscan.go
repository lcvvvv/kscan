package run

import (
	"fmt"
	"kscan/app"
	"kscan/lib/gonmap"
	"kscan/lib/pool"
	"kscan/lib/queue"
	"time"
)

type kscan struct {
	target *queue.Queue
	result *queue.Queue
	config app.Config
	pool   struct {
		host   *pool.Pool
		port   *pool.Pool
		banner *pool.Pool
	}
}

func New(config app.Config) *kscan {
	k := &kscan{
		target: queue.New(),
		config: config,
		result: queue.New(),
	}

	hostThreads := len(k.config.HostTarget)
	hostThreads = hostThreads/10 + 1
	if hostThreads > 400 {
		hostThreads = 400
	}

	k.pool.banner = pool.NewPool(config.Threads)
	k.pool.port = pool.NewPool(config.Threads)
	k.pool.host = pool.NewPool(hostThreads)
	return k
}

//func (k *kscan) Push(i interface{}) {
//	k.target.Push(i)
//}
//
//func (k *kscan) PushAll(iArr []string) {
//	for _, i := range iArr {
//		k.Push(i)
//	}
//}

//func (k *kscan) PushUrl(iArr []string) {
//	k.PushAll(k.PORT, iArr)
//}

func (k *kscan) HostDiscovery(hostArr []string, open bool) {
	k.pool.host.Function = func(i interface{}) interface{} {
		ip := i.(string)
		//如果关闭存活性检测，则默认所有IP存活
		if open == false {
			return ip
		}
		//经过存活性检测未存活的IP不会进行下一步测试
		if gonmap.HostDiscovery(ip) {
			return ip
		}
		return nil
	}

	//启用主机存活性探测任务下发器
	go func() {
		for _, host := range hostArr {
			k.pool.host.In <- host
		}
		//关闭主机存活性探测下发信道
		k.pool.host.InDone()
	}()
	//开始执行主机存活性探测任务
	k.pool.host.Run()
}

func (k *kscan) PortDiscovery() {
	k.pool.port.Function = func(i interface{}) interface{} {
		netloc := i.(string)
		if ok := gonmap.PortScan(netloc, time.Duration(k.config.Timeout)*time.Second); ok {
			return netloc
		}
		return nil
	}

	//启用端口存活性探测任务下发器
	go func() {
		for out := range k.pool.host.Out {
			for _, port := range k.config.Port {
				netloc := fmt.Sprintf("%s:%d", out, port)
				k.pool.port.In <- netloc
			}
		}
		k.pool.port.InDone()
	}()
	//开始执行端口存活性探测任务
	k.pool.port.Run()
}

func (k *kscan) GetPortBanner() {
	k.pool.banner.Function = func(i interface{}) interface{} {
		netloc := i.(string)
		r := GetPortBanner(netloc, gonmap.New())
		if r.Status == "OPEN" || r.Status == "MATCHED" {
			return r
		} else {
			return nil
		}
	}

	//banner识别任务下发终止器，结束信道
	isDone := make(chan bool)
	go func() {
		i := 0
		for range isDone {
			i++
			if i == 2 {
				break
			}
		}
		k.pool.banner.InDone()
	}()

	//指定Url任务下发器
	go func() {
		for _, url := range k.config.UrlTarget {
			k.pool.banner.In <- url
		}
		isDone <- true
	}()

	//启用协议识别任务下发器
	go func() {
		for out := range k.pool.port.Out {
			k.pool.banner.In <- out
		}
		isDone <- true
	}()

	//开始执行协议识别任务
	k.pool.banner.Run()
}
func (k *kscan) Output() {
	//输出协议识别结果
	for out := range k.pool.banner.Out {
		portinfo := out.(*PortInformation)
		portinfo.MakeInfo()
		fmt.Println(portinfo.Info)
	}
}

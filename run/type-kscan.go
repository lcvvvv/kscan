package run

import (
	"fmt"
	"github.com/lcvvvv/urlparse"
	"kscan/app"
	"kscan/lib/gonmap"
	"kscan/lib/pool"
	"kscan/lib/queue"
	"kscan/lib/slog"
)

type kscan struct {
	target *queue.Queue
	result *queue.Queue
	config app.Config
	pool   struct {
		host      *pool.Pool
		port      *pool.Pool
		tcpBanner *pool.Pool
		appBanner *pool.Pool
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

	k.pool.appBanner = pool.NewPool(config.Threads)
	k.pool.tcpBanner = pool.NewPool(config.Threads)
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
		slog.Info("主机存活性探测任务下发完毕，关闭信道")
		k.pool.host.InDone()
	}()
	//开始执行主机存活性探测任务
	k.pool.host.Run()
}

func (k *kscan) PortDiscovery() {
	k.pool.port.Function = func(i interface{}) interface{} {
		netloc := i.(string)
		if gonmap.PortScan(netloc, k.config.Timeout) {
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
		slog.Info("端口存活性探测任务下发完毕，关闭信道")
		k.pool.port.InDone()
	}()
	//开始执行端口存活性探测任务
	k.pool.port.Run()
}

func (k *kscan) GetTcpBanner() {
	k.pool.tcpBanner.Function = func(i interface{}) interface{} {
		netloc := i.(string)
		r := gonmap.GetTcpBanner(netloc, gonmap.New(), k.config.Timeout*10)
		return r
	}

	//启用TCP层面协议识别任务下发器
	go func() {
		for out := range k.pool.port.Out {
			k.pool.tcpBanner.In <- out
		}
		slog.Info("端口TCP层协议识别任务下发完毕，关闭信道")
		k.pool.tcpBanner.InDone()
	}()

	//开始执行TCP层面协议识别任务
	k.pool.tcpBanner.Run()
}

func (k *kscan) GetAppBanner() {
	k.pool.appBanner.Function = func(i interface{}) interface{} {
		var r *gonmap.AppBanner
		switch i.(type) {
		case string:
			url, _ := urlparse.Load(i.(string))
			r = gonmap.GetAppBannerFromUrl(url)
		case *gonmap.TcpBanner:
			tcpBanner := i.(*gonmap.TcpBanner)
			r = gonmap.GetAppBannerFromTcpBanner(tcpBanner)
		}
		return r
	}

	//appBanner识别任务下发终止器，结束信道
	isDone := make(chan bool)
	go func() {
		i := 0
		for range isDone {
			i++
			if i == 2 {
				break
			}
		}
		k.pool.appBanner.InDone()
		slog.Info("全部端口应用层协议识别任务下发完毕，关闭信道")
	}()

	//指定Url任务下发器
	go func() {
		for _, url := range k.config.UrlTarget {
			k.pool.appBanner.In <- url
		}
		isDone <- true
		slog.Info("自定义应用层协议识别任务下发完毕")
	}()

	//启用App层面协议识别任务下发器
	go func() {
		for out := range k.pool.tcpBanner.Out {
			k.pool.appBanner.In <- out
		}
		isDone <- true
		slog.Info("存活应用层协议识别任务下发完毕")
	}()

	//开始执行App层面协议识别任务
	k.pool.appBanner.Run()
}

func (k *kscan) Output() {
	//输出协议识别结果
	for out := range k.pool.appBanner.Out {
		a := out.(*gonmap.AppBanner)
		if a != nil {
			str := a.Output()
			slog.Data(str)
			if k.config.Output != nil {
				k.config.WriteLine(str)
			}
		}
	}
}

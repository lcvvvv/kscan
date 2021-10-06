package run

import (
	"encoding/json"
	"fmt"
	"kscan/app"
	"kscan/lib/IP"
	"kscan/lib/gonmap"
	"kscan/lib/hydra"
	"kscan/lib/misc"
	"kscan/lib/pool"
	"kscan/lib/queue"
	"kscan/lib/slog"
	"kscan/lib/urlparse"
	"os"
	"path"
	"strings"
	"sync"
	"time"
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
	watchDog struct {
		output  chan interface{}
		hydra   chan interface{}
		wg      *sync.WaitGroup
		trigger bool
	}
	hydra struct {
		pool  *pool.Pool
		queue *queue.Queue
		done  bool
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

	k.watchDog.hydra = make(chan interface{})
	k.watchDog.output = make(chan interface{})
	k.watchDog.wg = &sync.WaitGroup{}
	k.watchDog.trigger = false

	k.hydra.pool = pool.NewPool(10)
	k.hydra.queue = queue.New()
	k.hydra.done = false
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
		if open == true {
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
	slog.Info("主机存活性探测任务完成")
}

func (k *kscan) PortDiscovery() {
	k.pool.port.Function = func(i interface{}) interface{} {
		netloc := i.(string)
		protocol := "tcp"
		if port := strings.Split(netloc, ":")[1]; port == "161" {
			protocol = "udp"
		}
		if gonmap.PortScan(protocol, netloc, k.config.Timeout) {
			return netloc
		}
		return nil
	}

	//启用端口存活性探测任务下发器
	go func() {
		for out := range k.pool.host.Out {
			for _, port := range k.config.Port {
				if port == 161 {
					//如果是公网IP且使用默认端口扫描策略，则不会扫描161端口
					if IP.IsPrivateIPAddr(out.(string)) == false && len(app.Setting.Port) == 400 {
						continue
					}
				}
				netloc := fmt.Sprintf("%s:%d", out, port)
				k.pool.port.In <- netloc
			}
		}
		slog.Info("端口存活性探测任务下发完毕，关闭信道")
		k.pool.port.InDone()
	}()
	//开始执行端口存活性探测任务
	k.pool.port.Run()
	slog.Info("端口存活性探测任务完成")
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
		slog.Info("TCP层协议识别任务下发完毕，关闭信道")
		k.pool.tcpBanner.InDone()
	}()

	//开始执行TCP层面协议识别任务
	k.pool.tcpBanner.Run()
	slog.Info("TCP层协议识别任务完成")

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
			if tcpBanner == nil {
				return nil
			}
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
		slog.Info("应用层协议识别任务下发完毕，关闭信道")
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
		slog.Info("开放端口应用层协议识别任务下发完毕")
	}()

	//开始执行App层面协议识别任务
	k.pool.appBanner.Run()
	slog.Info("应用层协议识别任务完成")
}

func (k *kscan) GetAppBannerFromCheck() {
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

	//指定Url任务下发器
	go func() {
		for _, url := range k.config.UrlTarget {
			k.pool.appBanner.In <- url
		}
		k.pool.appBanner.InDone()
		slog.Info("应用层协议识别任务下发完毕，关闭信道")
	}()

	//开始执行App层面协议识别任务
	k.pool.appBanner.Run()
	slog.Info("应用层协议识别任务完成")
}

func (k *kscan) Output() {
	//输出协议识别结果
	var bannerMapArr []map[string]string
	for out := range k.watchDog.output {
		if out == nil {
			continue
		}
		var disp string
		var write string
		//打开触发器,若长时间无输出，触发器会输出进度
		k.watchDog.trigger = true
		//输出结果
		switch out.(type) {
		case *gonmap.AppBanner:
			banner := out.(*gonmap.AppBanner)
			if banner == nil {
				continue
			}
			bannerMapArr = append(bannerMapArr, banner.Map())
			write = banner.Output()
			disp = banner.Display()
		case hydra.AuthInfo:
			info := out.(hydra.AuthInfo)
			if info.Status == false {
				continue
			}
			write = info.Output()
			disp = info.Display()
		}
		slog.Data(disp)
		if k.config.Output != nil {
			k.config.WriteLine(write)
		}
	}
	//输出json
	if app.Setting.OutputJson != "" {
		fileName := app.Setting.OutputJson
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			_ = os.MkdirAll(path.Dir(fileName), os.ModePerm)
		}
		bytes, _ := json.Marshal(bannerMapArr)
		err := misc.WriteLine(fileName, bytes)
		if err == nil {
			slog.Infof("扫描完成，Json文件已输出至：", fileName)
		} else {
			slog.Warning("输出Json失败！错误信息：", err.Error())
		}
	}
}

func (k *kscan) WatchDog() {

	k.watchDog.wg.Add(1)
	//触发器校准，每隔60秒会将触发器关闭
	go func() {
		for true {
			time.Sleep(60 * time.Second)
			k.watchDog.trigger = false
		}
	}()
	//轮询触发器，每隔一段时间会检测触发器是否打开
	go func() {
		for true {
			time.Sleep(59 * time.Second)
			if k.watchDog.trigger == false {
				if num := k.pool.host.JobsList.Length(); num > 0 {
					i := k.pool.host.JobsList.Peek()
					info := i.(string)
					slog.Warningf("当前主机存活性检测任务未完成，其并发协程数为：%d，具体其中的一个协程信息为：%s", num, info)
					continue
				}
				if num := k.pool.port.JobsList.Length(); num > 0 {
					i := k.pool.port.JobsList.Peek()
					info := i.(string)
					slog.Warningf("当前端口存活性检测任务未完成，其并发协程数为：%d，具体其中的一个协程信息为：%s", num, info)
					continue
				}
				if num := k.pool.tcpBanner.JobsList.Length(); num > 0 {
					i := k.pool.tcpBanner.JobsList.Peek()
					info := i.(string)
					slog.Warningf("当前TCP层指纹识别任务未完成，其并发协程数为：%d，具体其中的一个协程信息为：%s", num, info)
					continue
				}
				if num := k.pool.appBanner.JobsList.Length(); num > 0 {
					i := k.pool.appBanner.JobsList.Peek()
					var info string
					switch i.(type) {
					case string:
						info = i.(string)
					case *gonmap.TcpBanner:
						tcpBanner := i.(*gonmap.TcpBanner)
						if tcpBanner == nil {
							continue
						}
						info = tcpBanner.Target.URI()
					}
					slog.Warningf("当前应用层指纹检测并发协程数为：%d，具体其中的一个协程信息为：%s", num, info)
					continue
				}
			}
		}
	}()
	//Hydra模块
	if app.Setting.Hydra {
		slog.Info("hydra模块已开启，开始监听暴力破解任务")
		k.watchDog.wg.Add(1)
	}

	for out := range k.pool.appBanner.Out {
		k.watchDog.output <- out
		if app.Setting.Hydra {
			k.watchDog.hydra <- out
		}
	}

	k.watchDog.wg.Done()
	close(k.watchDog.hydra)

	k.watchDog.wg.Wait()
	close(k.watchDog.output)
}

func (k *kscan) Hydra() {
	//初始化默认密码字典
	hydra.InitDefaultAuthMap()
	//加载自定义字典
	hydra.InitCustomAuthMap()
	//初始化变量
	k.hydra.pool.Function = func(i interface{}) interface{} {
		if i == nil {
			return nil
		}
		banner := i.(*gonmap.AppBanner)
		//适配爆破模块
		authInfo := hydra.NewAuthInfo(banner.IPAddr, banner.Port, banner.Protocol)
		crack := hydra.NewCracker(authInfo, 10)
		go crack.Run()
		//爆破结果获取
		var out hydra.AuthInfo
		for info := range crack.Out {
			out = info
		}
		return out
	}
	//暴力破解任务收集器
	go func() {
		for out := range k.watchDog.hydra {
			if out == nil {
				continue
			}
			banner := out.(*gonmap.AppBanner)
			if banner == nil {
				continue
			}
			if hydra.Ok(banner.Protocol, banner.Port) == false {
				continue
			}
			k.hydra.queue.Push(banner)
		}
		k.hydra.done = true
	}()
	//暴力破解任务下发器
	go func() {
		var TargetMap = make(map[string][]string)
		for true {
			if k.hydra.queue.Len() == 0 && k.hydra.done == true {
				break
			}
			pop := k.hydra.queue.Pop()
			if pop == nil {
				continue
			}
			banner := pop.(*gonmap.AppBanner)
			//若目标是第一次出现，则直接进行扫描
			if _, ok := TargetMap[banner.Netloc()]; ok == false {
				TargetMap[banner.Netloc()] = []string{banner.Protocol}
				k.hydra.pool.In <- banner
				continue
			}
			protocolArr := TargetMap[banner.Netloc()]
			if misc.IsInStrArr(protocolArr, banner.Protocol) == false {
				if arr := []string{"rdp", "smb"}; misc.IsInStrArr(arr, banner.Protocol) {
					TargetMap[banner.Netloc()] = append(protocolArr, arr...)
					k.hydra.pool.In <- banner
					continue
				}
				if arr := []string{"pop3", "smtp", "imap"}; misc.IsInStrArr(arr, banner.Protocol) {
					TargetMap[banner.Netloc()] = append(protocolArr, arr...)
					k.hydra.pool.In <- banner
					continue
				}
				TargetMap[banner.Netloc()] = append(protocolArr, banner.Protocol)
				k.hydra.pool.In <- banner
				continue
			}
		}
		//关闭输出信道
		k.hydra.pool.InDone()
	}()
	//暴力破解输出接受器
	go func() {
		for out := range k.hydra.pool.Out {
			k.watchDog.output <- out
		}
		k.watchDog.wg.Done()
	}()

	k.hydra.pool.Run()
}

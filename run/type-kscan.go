package run

import (
	"encoding/json"
	"fmt"
	"kscan/app"
	"kscan/lib/color"
	"kscan/lib/gonmap"
	"kscan/lib/httpfinger"
	"kscan/lib/hydra"
	"kscan/lib/misc"
	"kscan/lib/pool"
	"kscan/lib/queue"
	"kscan/lib/slog"
	"kscan/lib/smap"
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
		host struct {
			icmp *pool.Pool
			tcp  *pool.Pool
			Out  chan interface{}
		}
		port struct {
			tcp *pool.Pool
			Out chan interface{}
		}
		tcpBanner struct {
			tcp *pool.Pool
			Out chan interface{}
		}
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

	k.pool.tcpBanner.tcp = pool.NewPool(config.Threads)
	k.pool.tcpBanner.Out = make(chan interface{})

	k.pool.port.tcp = pool.NewPool(config.Threads)
	k.pool.port.Out = make(chan interface{})

	k.pool.host.icmp = pool.NewPool(hostThreads)
	k.pool.host.tcp = pool.NewPool(hostThreads)
	k.pool.host.Out = make(chan interface{})

	k.pool.appBanner.Interval = time.Microsecond * 500

	k.pool.tcpBanner.tcp.Interval = time.Microsecond * 500
	k.pool.port.tcp.Interval = time.Microsecond * 500

	k.watchDog.hydra = make(chan interface{})
	k.watchDog.output = make(chan interface{})
	k.watchDog.wg = &sync.WaitGroup{}
	k.watchDog.trigger = false

	k.hydra.pool = pool.NewPool(10)
	k.hydra.queue = queue.New()
	k.hydra.done = false
	return k
}

func (k *kscan) HostDiscovery(hostArr []string, open bool) {
	k.pool.host.icmp.Function = func(i interface{}) interface{} {
		ip := i.(string)
		host := NewHost(ip)
		//如果关闭存活性检测，则默认所有IP存活
		if open == true {
			return host.Up()
		}
		//经过存活性检测未存活的IP不会进行下一步测试
		if gonmap.HostDiscoveryForIcmp(ip) == true {
			slog.Debug(host.addr, " is alive")
			return host.Up()
		}
		//ICMP检测不存活的主机，将发送至TCP存活性检测
		k.pool.host.tcp.In <- ip
		return host.Down()
	}
	k.pool.host.tcp.Function = func(i interface{}) interface{} {
		ip := i.(string)
		host := NewHost(ip)
		//经过存活性检测未存活的IP不会进行下一步测试
		if gonmap.HostDiscoveryForTcp(ip) == true {
			slog.Debug(host.addr, " is alive")
			return host.Up()
		}
		return host.Down()
	}
	//开启主机存活性探测输出调度器
	go func() {
		//ICMP输出调度器
		wg := &sync.WaitGroup{}
		wg.Add(2)
		go func() {
			for out := range k.pool.host.icmp.Out {
				host := out.(*Host)
				if host.status == Up {
					k.pool.host.Out <- out
				}
			}
			wg.Done()
		}()
		//TCP输出调度器
		go func() {
			for out := range k.pool.host.tcp.Out {
				host := out.(*Host)
				if host.status == Up {
					k.pool.host.Out <- out
				}
			}
			wg.Done()
		}()
		wg.Wait()
		close(k.pool.host.Out)
	}()

	//启用ICMP主机存活性探测任务下发器
	go func() {
		for _, host := range hostArr {
			k.pool.host.icmp.In <- host
		}
		//关闭主机存活性探测下发信道
		if k.config.ClosePing == false {
			slog.Info("主机存活性探测任务下发完毕")
		}
		k.pool.host.icmp.InDone()
	}()

	//开始执行主机存活性探测任务
	k.pool.host.tcp.RunBack()
	k.pool.host.icmp.Run()
	k.pool.host.tcp.InDone()
	k.pool.host.tcp.Wait()
	if k.config.ClosePing == false {
		slog.Warning("主机存活性探测任务完成")
	}
}

func (k *kscan) PortDiscovery() {
	var upHosts = smap.New()
	//启用端口存活性探测任务下发器
	go func() {
		for out := range k.pool.host.Out {
			host := out.(*Host)
			upHosts.Set(host.addr, host)
			for _, port := range k.config.Port {
				//if port == 161 {
				//	//如果是公网IP且使用默认端口扫描策略，则不会扫描161端口
				//	if IP.IsPrivateIPAddr(out.(string)) == false && len(app.Setting.Port) == 400 {
				//		continue
				//	}
				//}
				netloc := NewPort(host.addr, port)
				k.pool.port.tcp.In <- netloc
			}
		}
		slog.Info("端口存活性探测任务下发完毕")
		k.pool.port.tcp.InDone()
	}()
	//启用端口存活性探测结果接受器
	go func() {
		for out := range k.pool.port.tcp.Out {
			netloc := out.(*Port)
			if value, ok := upHosts.Get(netloc.addr); ok {
				host := value.(*Host)
				host.SetPort(netloc.port, netloc.status)
				if netloc.status == Open {
					k.pool.port.Out <- netloc
					host.Up()
				}
				if host.IsOpenPort() == false && host.Length() == len(k.config.Port) && k.config.ClosePing == false {
					url := fmt.Sprintf("icmp://%s", host.addr)
					description := color.Red(color.Overturn("Not Open Any Port"))
					output := fmt.Sprintf("%-30v %-26v %s", url, "Up", description)
					k.watchDog.output <- output
				}
				upHosts.Set(host.addr, host)
			}
		}
		close(k.pool.port.Out)
	}()

	//定义端口存活性检测函数
	k.pool.port.tcp.Function = func(i interface{}) interface{} {
		netloc := i.(*Port)
		if netloc.port == 161 || netloc.port == 137 {
			return netloc.Unknown()
		}
		if gonmap.PortScan("tcp", netloc.UnParse(), k.config.Timeout) {
			slog.Debug(netloc, " is open")
			return netloc.Open()
		}
		return netloc.Close()
	}
	//开始执行端口存活性探测任务
	k.pool.port.tcp.Run()
	slog.Warning("端口存活性探测任务完成")
}

func (k *kscan) GetTcpBanner() {

	k.pool.tcpBanner.tcp.Function = func(i interface{}) interface{} {
		netloc := i.(*Port)
		r := gonmap.GetTcpBanner(netloc.UnParse(), gonmap.New(), k.config.Timeout*20)
		return r
	}

	//启用TCP层面协议识别任务下发器
	go func() {
		for out := range k.pool.port.Out {
			netloc := out.(*Port)
			if netloc.status == Close {
				continue
			}
			k.pool.tcpBanner.tcp.In <- out
		}
		slog.Info("TCP层协议识别任务下发完毕")
		k.pool.tcpBanner.tcp.InDone()
	}()

	//启用TCP层指纹探测结果接受器
	go func() {
		for out := range k.pool.tcpBanner.tcp.Out {
			if out == nil {
				continue
			}
			tcpBanner := out.(*gonmap.TcpBanner)
			if tcpBanner == nil {
				continue
			}

			uri := tcpBanner.Target.URI()
			status := tcpBanner.Status()
			service := tcpBanner.TcpFinger.Service
			slog.Debugf("%s %s %s", uri, status, service)
			k.pool.tcpBanner.Out <- out
		}
		close(k.pool.tcpBanner.Out)
	}()

	//开始执行TCP层面协议识别任务
	k.pool.tcpBanner.tcp.Run()
	slog.Warning("TCP层协议识别任务完成")

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
		slog.Info("应用层协议识别任务下发完毕")
	}()

	//指定Url任务下发器
	go func() {
		for _, url := range k.config.UrlTarget {
			k.pool.appBanner.In <- url
		}
		isDone <- true
		//slog.Info("自定义应用层协议识别任务下发完毕")
	}()

	//启用App层面协议识别任务下发器
	go func() {
		for out := range k.pool.tcpBanner.Out {
			k.pool.appBanner.In <- out
		}
		isDone <- true
		//slog.Info("开放端口应用层协议识别任务下发完毕")
	}()

	//开始执行App层面协议识别任务
	k.pool.appBanner.Run()
	slog.Warning("应用层协议识别任务完成")
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
		slog.Info("应用层协议识别任务下发完毕")
	}()

	//开始执行App层面协议识别任务
	k.pool.appBanner.Run()
	slog.Warning("应用层协议识别任务完成")
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
			write = banner.Output(app.Setting.CloseColor)
			disp = banner.Display(app.Setting.CloseColor)
		case hydra.AuthInfo:
			info := out.(hydra.AuthInfo)
			if info.Status == false {
				continue
			}
			write = info.Output()
			disp = info.Display()
		case string:
			outString := out.(string)
			write = outString
			disp = outString
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

	if len(httpfinger.NewKeywords) > 0 {
		newKeywords := misc.RemoveDuplicateElement(httpfinger.NewKeywords)
		slog.Warning("为了使kscan变得更好，请将finger.txt文件，提交到作者的Github")
		dir, _ := os.Getwd()
		slog.Warningf("发现新的http指纹[%d]条:%s/%s", len(newKeywords), dir, "finger.txt")
		data := strings.Join(newKeywords, "\r\n")
		_ = misc.WriteLine("finger.txt", []byte(data))
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
				//if num := k.pool.host.JobsList.Length(); num > 0 {
				//	i := k.pool.host.JobsList.Peek()
				//	info := i.(string)
				//	slog.Warningf("当前主机存活性检测任务未完成，其并发协程数为：%d，具体其中的一个协程信息为：%s", num, info)
				//	continue
				//}
				if num := k.pool.port.tcp.JobsList.Length(); num > 0 {
					i := k.pool.port.tcp.JobsList.Peek()
					info := i.(*Port)
					slog.Warningf("正在进行端口存活性检测，其并发协程数为：%d，具体其中的一个协程信息为：%s", num, info.UnParse())
					continue
				}
				if num := k.pool.tcpBanner.tcp.JobsList.Length(); num > 0 {
					i := k.pool.tcpBanner.tcp.JobsList.Peek()
					info := i.(*Port)
					slog.Warningf("正在进行TCP层指纹识别，其并发协程数为：%d，具体其中的一个协程信息为：%s", num, info)
					continue
				}
				if num := k.pool.appBanner.JobsList.Length(); num > 0 {
					i := k.pool.appBanner.JobsList.Peek()
					var info string
					switch i.(type) {
					case *Port:
						info = i.(*Port).UnParse()
					case *gonmap.TcpBanner:
						tcpBanner := i.(*gonmap.TcpBanner)
						if tcpBanner == nil {
							continue
						}
						info = tcpBanner.Target.URI()
					}
					slog.Warningf("正在进行应用层指纹识别，其并发协程数为：%d，具体其中的一个协程信息为：%s", num, info)
					continue
				}
			}
		}
	}()
	//Hydra模块
	if app.Setting.Hydra {
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
	slog.Info("hydra模块已开启，开始监听暴力破解任务")
	slog.Warning("当前已开启的hydra模块为：", misc.Intersection(hydra.ProtocolList, app.Setting.HydraMod))
	//初始化默认密码字典
	hydra.InitDefaultAuthMap()
	//加载自定义字典
	hydra.InitCustomAuthMap(app.Setting.HydraUser, app.Setting.HydraPass)
	//初始化变量
	k.hydra.pool.Function = func(i interface{}) interface{} {
		if i == nil {
			return nil
		}
		banner := i.(*gonmap.AppBanner)
		//适配爆破模块
		authInfo := hydra.NewAuthInfo(banner.IPAddr, banner.Port, banner.Protocol)
		crack := hydra.NewCracker(authInfo, app.Setting.HydraUpdate, 10)
		slog.Infof("[hydra]->开始对%v:%v[%v]进行暴力破解，字典长度为：%d", banner.IPAddr, banner.Port, banner.Protocol, crack.Length())
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
			if misc.IsInStrArr(app.Setting.HydraMod, banner.Protocol) == false {
				continue
			}
			if hydra.Ok(banner.Protocol) == false {
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

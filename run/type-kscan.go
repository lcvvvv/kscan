package run

import (
	"encoding/json"
	"fmt"
	"github.com/lcvvvv/gonmap"
	"github.com/lcvvvv/gonmap/lib/httpfinger"
	"github.com/lcvvvv/gonmap/lib/urlparse"
	"kscan/app"
	"kscan/core/hydra"
	"kscan/core/slog"
	"kscan/lib/color"
	"kscan/lib/misc"
	"kscan/lib/pool"
	"kscan/lib/queue"
	"kscan/lib/smap"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
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

	portScanMap *smap.SMap
}

func New(config app.Config) *kscan {
	k := &kscan{
		target: queue.New(),
		config: config,
		result: queue.New(),
	}

	hostThreads := len(k.config.HostTarget)
	hostThreads = hostThreads/5 + 1
	if hostThreads > 400 {
		hostThreads = 400
	}

	k.pool.appBanner = pool.NewPool(config.Threads)

	k.pool.tcpBanner.tcp = pool.NewPool(config.Threads)
	k.pool.tcpBanner.Out = make(chan interface{})

	k.pool.port.tcp = pool.NewPool(config.Threads * 4)
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

	k.portScanMap = smap.New()
	return k
}

func (k *kscan) HostDiscovery(hostArr []string, open bool) {
	k.pool.host.icmp.Function = func(i interface{}) interface{} {
		ip := i.(string)
		host := NewHost(ip, len(k.config.Port))
		//如果关闭存活性检测，则默认所有IP存活
		if open == true {
			return host.Up()
		}
		//经过存活性检测未存活的IP不会进行下一步测试
		if gonmap.HostDiscoveryForIcmp(ip) == true {
			slog.Println(slog.DEBUG, host.addr, " is alive")
			return host.Up()
		}
		//ICMP检测不存活的主机，将发送至TCP存活性检测
		k.pool.host.tcp.In <- ip
		return host.Down()
	}
	k.pool.host.tcp.Function = func(i interface{}) interface{} {
		ip := i.(string)
		host := NewHost(ip, len(k.config.Port))
		//经过存活性检测未存活的IP不会进行下一步测试
		if gonmap.HostDiscoveryForTcp(ip) == true {
			slog.Println(slog.DEBUG, host.addr, " is alive")
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
				if host.IsAlive() {
					k.pool.host.Out <- out
				}
			}
			wg.Done()
		}()
		//TCP输出调度器
		go func() {
			for out := range k.pool.host.tcp.Out {
				host := out.(*Host)
				if host.IsAlive() {
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
			slog.Println(slog.INFO, "主机存活性探测任务下发完毕")
		}
		k.pool.host.icmp.InDone()
	}()

	//开始执行主机存活性探测任务
	k.pool.host.tcp.RunBack()
	k.pool.host.icmp.Run()
	k.pool.host.tcp.InDone()
	k.pool.host.tcp.Wait()
	if k.config.ClosePing == false {
		slog.Println(slog.WARN, "主机存活性探测任务完成")
	}
}

func (k *kscan) PortDiscovery() {
	//启用端口存活性探测任务下发器
	go func() {
		var wg int32 = 0
		var threads = k.config.Threads

		for out := range k.pool.host.Out {
			host := out.(*Host)
			k.portScanMap.Set(host.addr, host)
			atomic.AddInt32(&wg, 1)
			go func() {
				defer func() { atomic.AddInt32(&wg, -1) }()
				for _, port := range k.config.Port {
					netloc := NewPort(host.addr, port)
					k.pool.port.tcp.In <- netloc
				}
			}()
			for int(wg) >= threads {
				time.Sleep(1 * time.Second)
			}
		}
		for wg > 0 {
			time.Sleep(1 * time.Second)
		}
		slog.Println(slog.INFO, "端口存活性探测任务下发完毕")
		k.pool.port.tcp.InDone()
	}()
	//启用端口存活性探测结果接受器
	go func() {
		for out := range k.pool.port.tcp.Out {
			port := out.(*Port)
			value, _ := k.portScanMap.Get(port.addr)
			host := value.(*Host)
			host.SetAlivePort(port.port, port.status)

			if port.status == Open {
				k.pool.port.Out <- port
				host.Up()
			}
			if port.status == Unknown {
				k.pool.port.Out <- port
			}
			if host.Map.Port.Length() == host.Length.Port {
				//所有端口检测完，表示该主机端口存活性检测已结束
				host.FinishPortScan()
				//输出没有开放任何端口的主机
				if host.IsOpenPort() == false && k.config.ClosePing == false {
					url := fmt.Sprintf("icmp://%s", host.addr)
					description := color.Red(color.Overturn("Not Open Any Port"))
					output := fmt.Sprintf("%-30v %-26v %s", url, "Up", description)
					k.watchDog.output <- output
				}
			}
			k.portScanMap.Set(port.addr, host)

		}
		close(k.pool.port.Out)
	}()
	//定义端口存活性检测函数
	k.pool.port.tcp.Function = func(i interface{}) interface{} {
		netloc := i.(*Port)
		if netloc.port == 161 || netloc.port == 137 {
			return netloc.Unknown()
		}
		if gonmap.PortScan("tcp", netloc.addr, netloc.port, 1*time.Second) {
			slog.Println(slog.DEBUG, netloc.UnParse(), " is open")
			return netloc.Open()
		}
		return netloc.Close()
	}
	//开始执行端口存活性探测任务
	k.pool.port.tcp.Run()
	slog.Println(slog.WARN, "端口存活性探测任务完成")
}

func (k *kscan) GetTcpBanner() {
	k.pool.tcpBanner.tcp.Function = func(i interface{}) interface{} {
		port := i.(*Port)
		var r = gonmap.NewTcpBanner(port.addr, port.port)
		//slog.Println(slog.DEBUG, port.UnParse(),port.Status())
		if port.status == Close {
			return r.CLOSED()
		}
		return gonmap.GetTcpBanner(port.addr, port.port, gonmap.New(), k.config.Timeout*20)
	}

	//启用TCP层面协议识别任务下发器
	go func() {
		for out := range k.pool.port.Out {
			k.pool.tcpBanner.tcp.In <- out
		}
		slog.Println(slog.INFO, "TCP层协议识别任务下发完毕")
		k.pool.tcpBanner.tcp.InDone()
	}()

	//启用TCP层指纹探测结果接受器
	go func() {
		for out := range k.pool.tcpBanner.tcp.Out {
			//应该不存在此返回项
			if out == nil {
				continue
			}
			tcpBanner := out.(*gonmap.TcpBanner)
			//应该不存在此返回项
			if tcpBanner == nil {
				continue
			}

			//此处开始，正式开始对输出结果进行处理

			port := tcpBanner.Target.Port()
			addr := tcpBanner.Target.Addr()

			value, _ := k.portScanMap.Get(addr)
			host := value.(*Host)

			if tcpBanner.Status() == gonmap.Matched {
				//slog.Printf(slog.DEBUG, "%s:%d %s %s", addr, port, status, service)
				k.pool.tcpBanner.Out <- tcpBanner
			} else {
				if (tcpBanner.Target.Port() == 161 || tcpBanner.Target.Port() == 137) && tcpBanner.Response.Length() == 0 {
					tcpBanner.CLOSED()
				}
				//slog.Println(slog.WARN, tcpBanner.Target.URI(), tcpBanner.StatusDisplay())
			}
			host.Map.Tcp.Set(port, tcpBanner)

			k.portScanMap.Set(addr, host)
			if host.PortScanIsFinish() == false {
				continue
			}
			if host.Map.Tcp.Length() == host.Length.Tcp && host.CountUnknownPorts() > 0 {
				host.status.tcpScan = true
				k.watchDog.output <- host.DisplayUnknownPorts()
			}
		}
		k.portScanMap.Range(
			func(key, value interface{}) bool {
				host := value.(*Host)
				if host.CountUnknownPorts() == 0 {
					return true
				}
				if host.status.tcpScan == false {
					k.watchDog.output <- host.DisplayUnknownPorts()
				}
				return true
			},
		)
		close(k.pool.tcpBanner.Out)
	}()

	//开始执行TCP层面协议识别任务
	k.pool.tcpBanner.tcp.Run()
	slog.Println(slog.WARN, "TCP层协议识别任务完成")

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
		slog.Println(slog.INFO, "应用层协议识别任务下发完毕")
	}()

	//指定Url任务下发器
	go func() {
		for _, url := range k.config.UrlTarget {
			k.pool.appBanner.In <- url
		}
		isDone <- true
	}()

	//启用App层面协议识别任务下发器
	go func() {
		for out := range k.pool.tcpBanner.Out {
			tcpBanner := out.(*gonmap.TcpBanner)
			if tcpBanner.Status() != gonmap.Matched {
				continue
			}
			k.pool.appBanner.In <- out
		}
		isDone <- true
	}()

	//开始执行App层面协议识别任务
	k.pool.appBanner.Run()
	slog.Println(slog.WARN, "应用层协议识别任务完成")
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
		slog.Println(slog.INFO, "应用层协议识别任务下发完毕")
	}()

	//开始执行App层面协议识别任务
	k.pool.appBanner.Run()
	slog.Println(slog.WARN, "应用层协议识别任务完成")
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
			if app.Setting.Match != "" && strings.Contains(banner.Response, app.Setting.Match) == false {
				continue
			}
			bannerMapArr = append(bannerMapArr, banner.Map())
			write = outputTcpBanner(banner, app.Setting.CloseColor)
			disp = displayTcpBanner(banner, app.Setting.CloseColor)
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
		slog.Println(slog.DATA, disp)
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
			slog.Printf(slog.INFO, "扫描完成，Json文件已输出至：", fileName)
		} else {
			slog.Println(slog.WARN, "输出Json失败！错误信息：", err.Error())
		}
	}

	if len(httpfinger.NewKeywords) > 0 {
		newKeywords := misc.RemoveDuplicateElement(httpfinger.NewKeywords)
		slog.Println(slog.WARN, "为了使kscan变得更好，请将finger.txt文件，提交到作者的Github")
		dir, _ := os.Getwd()
		slog.Printf(slog.WARN, "发现新的http指纹[%d]条:%s/%s", len(newKeywords), dir, "finger.txt")
		data := strings.Join(newKeywords, "\r\n")
		_ = misc.WriteLine("finger.txt", []byte(data))
	}

}

func (k *kscan) WatchDog() {
	k.watchDog.wg.Add(1)
	//触发器轮询时间
	waitTime := 30 * time.Second
	//轮询触发器，每隔一段时间会检测触发器是否打开
	go func() {
		for true {
			time.Sleep(waitTime)
			if k.watchDog.trigger == false {
				slog.Printf(slog.WARN,
					"当前运行情况为:主机存活性检测并发【%d】个,端口存活性检测并发【%d】个,TCP层检测并发【%d】个,APP层检测并发【%d】个",
					k.pool.host.icmp.JobsList.Length()+k.pool.host.tcp.JobsList.Length(),
					k.pool.port.tcp.JobsList.Length(),
					k.pool.tcpBanner.tcp.JobsList.Length(),
					k.pool.appBanner.JobsList.Length(),
				)
			}
		}
	}()
	time.Sleep(time.Millisecond * 500)
	//触发器校准，每隔一段时间将触发器关闭
	go func() {
		for true {
			time.Sleep(waitTime)
			k.watchDog.trigger = false
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
	slog.Println(slog.INFO, "hydra模块已开启，开始监听暴力破解任务")
	slog.Println(slog.WARN, "当前已开启的hydra模块为：", misc.Intersection(hydra.ProtocolList, app.Setting.HydraMod))
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
		slog.Printf(slog.INFO, "[hydra]->开始对%v:%v[%v]进行暴力破解，字典长度为：%d", banner.IPAddr, banner.Port, banner.Protocol, crack.Length())
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

func displayTcpBanner(appBanner *gonmap.AppBanner, keyPrint bool) string {
	m := misc.FixMap(appBanner.FingerPrint())
	fingerPrint := color.StrMapRandomColor(m, keyPrint, []string{"ProductName", "Hostname", "DeviceType"}, []string{"ApplicationComponent"})
	fingerPrint = misc.FixLine(fingerPrint)
	format := "%-30v %-" + strconv.Itoa(misc.AutoWidth(appBanner.AppDigest, 26)) + "v %s"
	s := fmt.Sprintf(format, appBanner.URL(), appBanner.AppDigest, fingerPrint)
	return s
}

func outputTcpBanner(appBanner *gonmap.AppBanner, keyPrint bool) string {
	fingerPrint := misc.StrMap2Str(appBanner.FingerPrint(), keyPrint)
	fingerPrint = misc.FixLine(fingerPrint)
	s := fmt.Sprintf("%s\t%d\t%s\t%s", appBanner.URL(), appBanner.StatusCode, appBanner.AppDigest, fingerPrint)
	return s
}

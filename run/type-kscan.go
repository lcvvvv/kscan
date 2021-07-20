package run

import (
	"encoding/json"
	"fmt"
	"github.com/lcvvvv/urlparse"
	"kscan/app"
	"kscan/lib/IP"
	"kscan/lib/gonmap"
	"kscan/lib/hydra"
	"kscan/lib/misc"
	"kscan/lib/pool"
	"kscan/lib/queue"
	"kscan/lib/slog"
	"os"
	"path"
	"strings"
	"sync"
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
		output chan interface{}
		hydra  chan interface{}
		wg     *sync.WaitGroup
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
		slog.Info("自定义应用层协议识别任务下发完毕")
	}()

	//开始执行App层面协议识别任务
	k.pool.appBanner.Run()
}

func (k *kscan) Output() {
	//输出协议识别结果
	var bannerMapArr []map[string]string
	for out := range k.watchDog.output {
		if out == nil {
			continue
		}
		var disp string
		switch out.(type) {
		case *gonmap.AppBanner:
			banner := out.(*gonmap.AppBanner)
			if banner == nil {
				return
			}
			bannerMapArr = append(bannerMapArr, banner.Map())
			disp = banner.Output()
		case hydra.AuthInfo:
			info := out.(hydra.AuthInfo)
			disp = info.Output()
		}
		slog.Data(disp)
		if k.config.Output != nil {
			k.config.WriteLine(disp)
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

	if app.Setting.Hydra {
		slog.Info("hydra模块已开启，开始监听暴力破解任务")
		k.watchDog.wg.Add(1)
	}

	go func() {
		for out := range k.pool.appBanner.Out {
			k.watchDog.output <- out
			if app.Setting.Hydra {
				k.watchDog.hydra <- out
			}
		}
		k.watchDog.wg.Done()
		close(k.watchDog.hydra)
	}()
	k.watchDog.wg.Wait()
	close(k.watchDog.output)
}

func (k *kscan) Hydra() {
	//初始化默认密码字典
	hydra.InitDefaultAuthMap()
	//开始监听暴力破解任务
	for out := range k.watchDog.hydra {
		banner := out.(*gonmap.AppBanner)
		if banner == nil {
			continue
		}
		if hydra.Ok(banner.Protocol, banner.Port) == false {
			continue
		}
		//适配爆破模块
		authInfo := hydra.NewAuthInfo(banner.IPAddr, banner.Port, banner.Protocol)
		crack := hydra.NewCracker(authInfo, 10)
		go crack.Run()
		//爆破结果获取
		for info := range crack.Out {
			k.watchDog.output <- info
		}
	}
	k.watchDog.wg.Done()
}

package run

import (
	"kscan/app"
	"kscan/lib/slog"
	"time"
)

func Start(config app.Config) {
	K := New(config)

	//STEP0: 扫描结果调度器
	time.Sleep(time.Microsecond * 200)
	go K.WatchDog()

	if config.Check {
		slog.Warning("当前为验证模式，不会进行端口扫描，仅对给定地址进行指纹识别")
		//STEP1:验证模式直接进行应用层指纹识别
		time.Sleep(time.Microsecond * 200)
		go K.GetAppBannerFromCheck()
	} else {
		//STEP1:主机存活性检测
		time.Sleep(time.Microsecond * 200)
		go K.HostDiscovery(K.config.HostTarget, config.ClosePing)

		//STEP2：端口存活性检测
		time.Sleep(time.Microsecond * 200)
		go K.PortDiscovery()

		//STEP3：TCP指纹识别
		time.Sleep(time.Microsecond * 200)
		go K.GetTcpBanner()

		//STEP3: 应用层指纹识别
		time.Sleep(time.Microsecond * 200)
		go K.GetAppBanner()
	}

	//STEP_A: 暴力破解
	time.Sleep(time.Microsecond * 200)
	if app.Setting.Hydra {
		time.Sleep(time.Microsecond * 200)
		go K.Hydra()
	}
	//STEP_B: 输出
	K.Output()
}

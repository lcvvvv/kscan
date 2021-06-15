package run

import (
	"kscan/app"
	"kscan/lib/slog"
	"time"
)

func Start(config app.Config) {
	K := New(config)

	if config.Check {
		slog.Infof("当前为验证模式，不会进行端口扫描，仅对给定URL地址进行指纹识别")
		time.Sleep(time.Microsecond * 200)
		go K.GetAppBannerFromCheck()
		////STEP4: 输出
		K.Output()
		return
	}

	//STEP1:主机存活性检测
	time.Sleep(time.Microsecond * 200)
	go K.HostDiscovery(K.config.HostTarget, config.ScanPing)

	//STEP2：端口存活性检测
	time.Sleep(time.Microsecond * 200)
	go K.PortDiscovery()

	//STEP3：TCP指纹识别
	time.Sleep(time.Microsecond * 200)
	go K.GetTcpBanner()

	//STEP3: 应用层指纹识别
	time.Sleep(time.Microsecond * 200)
	go K.GetAppBanner()
	////STEP4: 输出
	K.Output()
}

package run

import (
	"kscan/app"
	"time"
)

func Start(config app.Config) {
	K := New(config)
	//STEP1:主机存活性检测
	time.Sleep(time.Microsecond * 200)
	go K.HostDiscovery(K.config.HostTarget, true)

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

package spy

import (
	"fmt"
	"kscan/lib/IP"
	"kscan/lib/gonmap"
	"kscan/lib/misc"
	"kscan/lib/pool"
	"kscan/lib/slog"
	"net"
	"strings"
)

var (
	Keyword = ""
	Scan    = false
	Target  []string
)

func Start() {
	slog.Info("将开始自动化存活网段探测，请注意，该模式会发送大量数据包，极易被风控感知，请慎用")
	slog.Info("现在开始进行自动网络环境探测")
	internet := internetTesting()
	_ = dnsTesting()
	var gatewayArr, All []string
	//若spy参数格式为IP地址，则将对指定的IP地址进行B段存活网关探测
	if IP.FormatCheck(Keyword) {
		slog.Infof("现在开始指定网段：%s，B段存活网关探测", Keyword)
		gatewayArr = IP.GetGatewayList(Keyword, "b")
		HostDiscoveryIcmpPool(gatewayArr)
		return
	}
	//依据情况判断是否进行172B段存活网关探测
	if Keyword == "all" || Keyword == "172" {
		//探测172段，B段存活网关
		slog.Infof("当前spy参数值为%s，将开始172段，大B段存活网关探测，此探测时间较长，请耐心等待", Keyword)
		gatewayArr = []string{}
		for i := 16; i <= 31; i++ {
			slog.Infof("现在开始枚举常见网段172.%d.0.0", i)
			gatewayArr = IP.GetGatewayList(fmt.Sprintf("172.%d.0.0", i), "b")
			gatewayArr = misc.RemoveDuplicateElementForMultiple(gatewayArr, All)
			if len(gatewayArr) > 0 {
				HostDiscoveryIcmpPool(gatewayArr)
				All = append(All, gatewayArr...)
			} else {
				slog.Info("该网段在之前已经枚举，此处不将不再重复枚举")
			}
		}
	}
	//依据情况判断是否进行10A段存活网关探测
	if Keyword == "all" || Keyword == "10" {
		//探测10段，A段存活网关
		slog.Infof("当前spy参数值为%s，将开始10段，A段存活网关探测，此探测时间较长，请耐心等待", Keyword)
		slog.Info("现在开始枚举常见网段10.0.0.0")
		gatewayArr = IP.GetGatewayList("10.0.0.1", "a")
		HostDiscoveryIcmpPool(gatewayArr)
	}
	//依据情况判断是否进行常规探测
	if Keyword == "all" || Keyword == "" {
		//探测网卡所在网段
		slog.Info("现在开始当前所在网段的B段网关存活性探测")
		gatewayArr = makeInterfaceGatwayList()
		gatewayArr = misc.RemoveDuplicateElement(gatewayArr)
		//探测当前所在网段B段网关
		gatewayArr = misc.RemoveDuplicateElementForMultiple(gatewayArr, All)
		if len(gatewayArr) > 0 {
			HostDiscoveryIcmpPool(gatewayArr)
			All = append(All, gatewayArr...)
		} else {
			slog.Info("该网段在之前已经枚举，此处不将不再重复枚举")
		}
		//探测存在特殊规律的网段
		if internet == false {
			slog.Info("现在开始枚举特殊网段1.1.1.0-255.255.255.0")
			gatewayArr = append(IP.GetGatewayList("1.1.1.1", "s"))
			HostDiscoveryIcmpPool(gatewayArr)
		}
	}

	//依据情况判断是否进行192B段存活网关探测
	if Keyword == "all" || Keyword == "" || Keyword == "192" {
		//探测常见网段192段，B段存活网关
		slog.Info("现在开始枚举常见网段192.168.0.0")
		gatewayArr = IP.GetGatewayList("192.168.0.1", "b")
		gatewayArr = misc.RemoveDuplicateElementForMultiple(gatewayArr, All)
		if len(gatewayArr) > 0 {
			HostDiscoveryIcmpPool(gatewayArr)
			All = append(All, gatewayArr...)
		} else {
			slog.Info("该网段在之前已经枚举，此处不将不再重复枚举")
		}
	}
	slog.Info("自动化存活网段探测结束")
	slog.Info("小提示：若需要对指定某一b段探测，可设置spy参数值为该段任意IP地址")
}

func makeInterfaceGatwayList() []string {
	var gatewayArr []string
	up, down := getInterfaces()
	for _, ip := range up {
		if strings.Contains(ip, "169.254") {
			continue
		}
		gatewayArr = append(gatewayArr, IP.GetGatewayList(ip, "b")...)
	}
	for _, ip := range down {
		if strings.Contains(ip, "169.254") {
			continue
		}
		gatewayArr = append(gatewayArr, IP.GetGatewayList(ip, "b")...)
	}
	return gatewayArr
}

func internetTesting() bool {
	if gonmap.HostDiscoveryForIcmp("114.114.114.114") {
		slog.Data("Internet--------[√]")
		return true
	} else {
		slog.Data("Internet--------[×]")
		return false
	}
}

func dnsTesting() bool {
	_, err := net.ResolveIPAddr("ip", "www.baidu.com")
	if err != nil {
		slog.Data("DNS-------------[×]")
		return false
	}
	slog.Data("DNS-------------[√]")
	return true
}

func getInterfaces() (up []string, down []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return up, down
	}
	for i := 0; i < len(netInterfaces); i++ {
		addrs, _ := netInterfaces[i].Addrs()
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					if (netInterfaces[i].Flags & net.FlagUp) != 0 {
						up = misc.UniStrAppend(up, ipnet.IP.String())
					} else {
						down = misc.UniStrAppend(down, ipnet.IP.String())
					}
				}
			}
		}
	}
	return up, down
}

func HostDiscoveryIcmpPool(gatewayArr []string) {
	spyPool := pool.NewPool(200)
	spyPool.Function = func(i interface{}) interface{} {
		ip := i.(string)
		//经过存活性检测未存活的IP不会进行下一步测试
		if gonmap.HostDiscoveryForIcmp(ip) {
			return ip
		}
		return nil
	}
	//启用ICMP存活性探测任务下发器
	go func() {
		for _, ip := range gatewayArr {
			spyPool.In <- ip
		}
		//关闭ICMP存活性探测下发信道
		spyPool.InDone()
	}()
	//开始执行主机存活性探测任务
	go spyPool.Run()
	//开始监测输出结果
	for out := range spyPool.Out {
		if out == nil {
			continue
		}
		ip := out.(string)
		slog.Data(ip)
		if Scan == false {
			continue
		}
		ipArr := IP.ExprToList(fmt.Sprintf("%s/24", ip))
		Target = append(Target, ipArr...)
	}
	if Scan == true {
		Target = misc.RemoveDuplicateElement(Target)
	}

}

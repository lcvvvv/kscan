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

func Start() {
	slog.Info("将开始自动化存活网段探测，请注意，该模式会发送大量数据包，极易被风控感知，请慎用")
	slog.Info("现在开始进行自动网络环境探测")
	internet := internetTesting()
	_ = dnsTesting()
	up, down := getInterfaces()
	var gatewayArr []string
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
	gatewayArr = misc.RemoveDuplicateElement(gatewayArr)

	slog.Info("现在开始当前所在网段的B段网关存活性探测")
	HostDiscoveryIcmpPool(gatewayArr)

	slog.Info("现在开始枚举常见网段192.168.0.0")
	gatewayArr = IP.GetGatewayList("192.168.0.1", "b")
	HostDiscoveryIcmpPool(gatewayArr)

	slog.Info("现在开始枚举常见网段10.0.0.0")
	gatewayArr = IP.GetGatewayList("10.0.0.1", "a")
	HostDiscoveryIcmpPool(gatewayArr)

	slog.Info("现在开始枚举常见网段172.16.0.0-172.31.0.0")
	gatewayArr = []string{}
	for i := 16; i <= 31; i++ {
		gatewayArr = append(IP.GetGatewayList(fmt.Sprintf("172.%d.0.0", i), "b"))
	}
	HostDiscoveryIcmpPool(gatewayArr)

	if internet == false {
		slog.Info("现在开始枚举特殊网段1.1.1.0-255.255.255.0")
		gatewayArr = append(IP.GetGatewayList("1.1.1.1", "s"))
		HostDiscoveryIcmpPool(gatewayArr)
	}

	slog.Info("自动化存活网段探测结束")
}

func internetTesting() bool {
	if gonmap.HostDiscovery("114.114.114.114") {
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
		if gonmap.HostDiscoveryIcmp(ip) {
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
		if out != nil {
			ip := out.(string)
			slog.Data(ip)
		}
	}
}

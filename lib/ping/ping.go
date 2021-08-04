package ping

import (
	"github.com/go-ping/ping"
	"kscan/lib/slog"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var linuxInit = false

func Check(ip string) bool {
	p, err := ping.NewPinger(ip)
	if err != nil {
		slog.Debug(err.Error())
		return false
	}
	if runtime.GOOS == "windows" {
		p.SetPrivileged(true)
	}
	if runtime.GOOS == "linux" && linuxInit == false {
		cmd := exec.Command("sysctl", "-n", "net.ipv4.ping_group_range")
		buf, _ := cmd.Output() // 错误处理略
		str := string(buf)
		if strings.Contains(str, "2147483647") == true {
			return true
		}
		cmd = exec.Command("cat", "/proc/sys/net/ipv4/ping_group_range")
		buf, _ = cmd.Output() // 错误处理略
		str = string(buf)
		if strings.Contains(str, "2147483647") == true {
			return true
		}
		linuxInit = false
	}
	p.Count = 2
	p.Timeout = time.Second * 2
	err = p.Run() // Blocks until finished.
	if err != nil {
		slog.Debug(err.Error())
	}
	s := p.Statistics()
	if s.PacketsRecv > 0 {
		return true
	}
	return false
}

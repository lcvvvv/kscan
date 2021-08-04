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
		linuxInit = func() bool {
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
			slog.Error("检测到当前操作系统为linux，Ping模块运行不正常，请执行下列命令后，再次打开此程序:\n",
				"sudo sysctl -w net.ipv4.ping_group_range=\"0 2147483647\"",
			)
			return false
		}()
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

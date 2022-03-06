package gosping

import (
	"os/exec"
	"runtime"
)

func OsPing(host string) bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", host, "-n", "1", "-w", "200")
	case "linux":
		cmd = exec.Command("ping", host, "-c", "1", "-w", "200", "-W", "200")
	case "darwin":
		cmd = exec.Command("ping", host, "-c", "1", "-W", "200")
	default:
		cmd = exec.Command("ping", host, "-c", "1", "-w", "200", "-W", "200")
	}
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

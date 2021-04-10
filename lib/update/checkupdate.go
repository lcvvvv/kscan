package update

import (
	"bufio"
	"kscan/lib/misc"
	"kscan/lib/slog"
	"os"
	"strings"
)

func CheckUpdate() {
	checkHeaderkeys()
}

func checkHeaderkeys() {
	if misc.FileIsExist("newHeaderkeys.txt") {
		reader := bufio.NewReader(os.Stdin)
		yesArr := []string{"", "\r", "yes", "y"}
		noArr := []string{"no", "n"}
		for {
			slog.Warning("检测到存在newHeaderkeys.txt文件，这里面存储了一些我们所不知道的HttpResponse头部参数和值，是否愿意发送给我们以便于我们更新更强大的规则库?[Yes/no]")
			text, _ := reader.ReadString('\n')
			text = misc.FixLine(text)
			text = strings.ToLower(text)
			if misc.IsInStrArr(yesArr, text) {
				slog.Info("正在上传规则库，感谢您的支持...")
				slog.Warning("成功上传规则库，正在删除newHeaderkeys.txt文件...")
				_ = os.Remove("newHeaderkeys.txt")
				break
			}
			if misc.IsInStrArr(noArr, text) {
				slog.Debug("收到，已为您删除newHeaderkeys.txt，将不会进行上传...")
				_ = os.Remove("newHeaderkeys.txt")
				break
			}
		}
	}
}

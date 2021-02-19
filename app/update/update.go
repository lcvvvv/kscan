package update

import (
	"../../lib/misc"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func CheckUpdate() {
	checkHeaderkeys()
}

func checkHeaderkeys() {
	if misc.FileIsExist("newHeaderkeys.txt") {
		reader := bufio.NewReader(os.Stdin)
		yesArr := []string{"", "yes", "y"}
		noArr := []string{"no", "n"}
		for {
			fmt.Print("[*]检测到存在newHeaderkeys.txt文件，这里面存储了一些我们所不知道的HttpResponse头部参数和值，是否愿意发送给我们以便于我们更新更强大的规则库?[Yes/no]")
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			text = strings.ToLower(text)
			if misc.IsInStrArr(yesArr, text) {
				fmt.Print("[+]正在上传规则库，感谢您的支持...\n")
				fmt.Print("[*]成功上传规则库，正在删除newHeaderkeys.txt文件...\n")
				_ = os.Remove("newHeaderkeys.txt")
				break
			}
			if misc.IsInStrArr(noArr, text) {
				fmt.Print("[-]收到，已为您删除newHeaderkeys.txt，将不会进行上传...\n")
				_ = os.Remove("newHeaderkeys.txt")
				break
			}
		}
	}
}

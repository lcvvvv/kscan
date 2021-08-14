package chinese

import (
	"fmt"
	"testing"
)

func TestRune(t *testing.T) {
	for i := 0; i <= 65535; i++ {
		srcRune := rune(i)
		fmt.Println(srcRune)
		fmt.Println(string(srcRune))
	}

	//s := "测试"
	//// 将字符串转换为rune数组
	//srcRunes := []rune(s)
	//fmt.Println(srcRunes)
	//// 创建一个新的rune数组，用来存放过滤后的数据
	//dstRunes := make([]rune, 0, len(srcRunes))
	//// 过滤不可见字符，根据上面的表的0-32和127都是不可见的字符
	//for _, c := range srcRunes {
	//	if c >= 0 && c <= 31 {
	//		continue
	//	}
	//	if c == 127 {
	//		continue
	//	}
	//	dstRunes = append(dstRunes, c)
	//}
	//fmt.Println(string(dstRunes))
}

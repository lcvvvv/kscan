package column

import (
	"fmt"
	"kscan/lib/misc"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	number := "1"
	char := "a"
	chinese := "汉"

	strArr := []string{number, char, chinese}

	for _, i := range misc.Xrange(10) {
		for _, v := range strArr {
			str := strings.Repeat(v, i)
			fmt.Println(str, "\t|")
			fmt.Println(len(str))
			fmt.Println(len([]rune(str)))
		}
	}
}

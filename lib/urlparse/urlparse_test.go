package urlparse

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	var s string
	s = "www.baidu.com:443"
	a, err := Load(s)
	fmt.Println(err)
	fmt.Printf("\n%#v", a)
	a.Scheme = ""
	fmt.Printf("\n%#v\n", a)
	fmt.Println(a.UnParse())

}

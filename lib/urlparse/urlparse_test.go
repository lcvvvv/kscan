package urlparse

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	var s string
	s = "114.114.114.114:144"
	a, err := Load(s)
	fmt.Println(err)
	fmt.Printf("\n%#v", a)
	a.Scheme = ""
	fmt.Printf("\n%#v\n", a)
	fmt.Println(a.UnParse())

}

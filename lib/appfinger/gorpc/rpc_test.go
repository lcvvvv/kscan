package gorpc

import (
	"fmt"
	"testing"
)

func TestNetBIOS(t *testing.T) {
	s, err := GetHostname("51.79.17.189")
	if s == nil {
		fmt.Println(err)
	}
	fmt.Println(s)
	return
}

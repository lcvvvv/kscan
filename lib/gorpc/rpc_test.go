package gorpc

import (
	"fmt"
	"testing"
)

func TestNetBIOS(t *testing.T) {
	finger, err := GetFinger("51.79.17.189")
	if finger == nil {
		fmt.Println(err)
	}
	fmt.Println(finger.Value())
	return
}

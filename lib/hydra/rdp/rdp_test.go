package rdp

import (
	"fmt"
	"testing"
)

func TestRdp(t *testing.T) {
	r, err := Check("192.168.217.147", "", "Administrator", "zaq1@WSX", 3389)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
}

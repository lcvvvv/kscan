package telnet

import (
	"fmt"
	"testing"
)

func TestTelnet(t *testing.T) {
	c := New("220.180.208.144", 23)
	err := c.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	buf, err := c.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(buf))
}

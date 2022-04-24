package qqwry

import (
	"fmt"
	"testing"
)

func TestQQwry(t *testing.T) {
	qqWry, err := NewQQwry("./qqwry.dat")
	if err != nil {
		panic(err)
	}
	result, err := qqWry.Find("114.114.114.114")
	if err != nil {
		panic(err)
	}
	fmt.Println(result.String())
}

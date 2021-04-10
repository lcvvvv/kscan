package main

import (
	"fmt"
	"kscan/lib/httpfinger"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println(httpfinger.FaviconHash)
}

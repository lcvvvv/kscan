package oracle

import (
	"fmt"
	"testing"
)

func TestGetSID(t *testing.T) {
	target := []string{""}
	for _, ip := range target {
		sid := GetSID(ip, 1521, ServiceName)
		if sid != "" {
			fmt.Println(ip, "\t", sid)
		}
	}

}

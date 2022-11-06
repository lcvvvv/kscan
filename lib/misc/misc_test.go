package misc

import (
	"fmt"
	"testing"
)

func TestFixLine(t *testing.T) {

	var s = "1 1  1          1"
	fmt.Println(s)
	fmt.Println(FixLine(s))

}

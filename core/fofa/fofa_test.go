package fofa

import (
	"fmt"
	"kscan/lib/misc"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGetPortMap(t *testing.T) {
	email := os.Getenv("FOFA_EMAIL")
	key := os.Getenv("FOFA_KEY")
	f := New(email, key)
	var fileSlice []string
	for i := 1; i <= 65535; i++ {
		keyword := "port=" + strconv.Itoa(i)
		responseJson := f.Search(keyword)
		row := fmt.Sprintf("%d\t%d", i, responseJson.Size)
		fmt.Println(row)
		fileSlice = append(fileSlice, row)
	}
	_ = misc.WriteLine("fofa_port.txt", []byte(strings.Join(fileSlice, "\n")))
}

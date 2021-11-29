package touch

import (
	"fmt"
	"kscan/lib/gonmap/simplenet"
	"strconv"
	"time"
)

type Response struct {
	Status bool
	Text   string
	Length int
}

func Touch(netloc string) Response {
	response, err := simplenet.Send("tcp", netloc, "", time.Second*3, 2048)
	if err != nil {
		fmt.Println(err)
		return Response{false, err.Error(), 0}
	}
	responseBuf := []byte(response)
	printStr := ""
	for _, charBuf := range responseBuf {
		if strconv.IsPrint(rune(charBuf)) {
			if charBuf > 0x7f {
				printStr += "?"
			} else {
				printStr += string(charBuf)
			}
			continue
		}
		printStr += fmt.Sprintf("\\x%x", string(charBuf))
	}
	return Response{true, printStr, len(response)}

}

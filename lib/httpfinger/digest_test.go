package httpfinger

import (
	"fmt"
	"testing"
)

func TestDigest(t *testing.T) {
	header := `Date: Wed, 09 Feb 2022 08:34:55 GMT
Server: Apache/2.4.7 (Ubuntu)
Cache-Control: private
Vary: Accept-Encoding
Content-Type: text/html; charset=utf-8
X-Powered-By: ThinkPHP
Expires: Thu, 19 Nov 1981 08:52:00 GMT
Set-Cookie: PHPSESSID=1hsk6e1vbm82bptq7f0gq2r9t2; path=/
Pragma: no-cache
`
	if xPoweredByRegx.MatchString(header) {
		server := xPoweredByRegx.FindStringSubmatch(header)[1]
		fmt.Println(server)
	}

	fmt.Println("abcdefg"[:1])
	fmt.Println("abcdefg"[1:])
}

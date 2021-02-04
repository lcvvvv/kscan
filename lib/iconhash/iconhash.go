package iconhash

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/twmb/murmur3"
	"io"
	"io/ioutil"
)

var isUint32 bool

func standBase64(braw []byte) []byte {
	bckd := base64.StdEncoding.EncodeToString(braw)
	var buffer bytes.Buffer
	for i := 0; i < len(bckd); i++ {
		ch := bckd[i]
		buffer.WriteByte(ch)
		if (i+1)%76 == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')
	return buffer.Bytes()
}

func mmh3Hash32(raw []byte) string {
	var h32 = murmur3.New32()
	_, _ = h32.Write(raw)
	if isUint32 {
		return fmt.Sprintf("%d", h32.Sum32())
	}
	return fmt.Sprintf("%d", int32(h32.Sum32()))
}

func Get(r io.Reader) (string, error) {
	body, err := ioutil.ReadAll(r)
	return mmh3Hash32(standBase64(body)), err
}

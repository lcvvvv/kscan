package chinese

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func ByteToGBK(strBuf []byte) []byte {
	if isUtf8(strBuf) {
		if GBKBuf, err := simplifiedchinese.GBK.NewEncoder().Bytes(strBuf); err == nil {
			return GBKBuf
		}
		if GB18030Buf, err := simplifiedchinese.GB18030.NewEncoder().Bytes(strBuf); err == nil {
			return GB18030Buf
		}
		if HZGB2312Buf, err := simplifiedchinese.HZGB2312.NewEncoder().Bytes(strBuf); err == nil {
			return HZGB2312Buf
		}
		return strBuf
	} else {
		return strBuf
	}
}

func ByteToUTF8(strBuf []byte) []byte {
	if isUtf8(strBuf) {
		return strBuf
	} else {
		if GBKBuf, err := simplifiedchinese.GBK.NewDecoder().Bytes(strBuf); err == nil {
			return GBKBuf
		}
		if GB18030Buf, err := simplifiedchinese.GB18030.NewDecoder().Bytes(strBuf); err == nil {
			return GB18030Buf
		}
		if HZGB2312Buf, err := simplifiedchinese.HZGB2312.NewDecoder().Bytes(strBuf); err == nil {
			return HZGB2312Buf
		}
		return strBuf
	}
}

func ToGBK(str string) string {
	strBuf := []byte(str)
	GBKBuf := ByteToGBK(strBuf)
	return string(GBKBuf)

}

func ToUTF8(str string) string {
	strBuf := []byte(str)
	Utf8Buf := ByteToUTF8(strBuf)
	return string(Utf8Buf)
}

func isUtf8(data []byte) bool {
	for i := 0; i < len(data); {
		if data[i]&0x80 == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if data[i]&0xc0 != 0x80 {
					return false
				}
				i++
			}
		} else {
			//其他情况说明不是utf-8
			return false
		}
	}
	return true
}

func preNUm(data byte) int {
	str := fmt.Sprintf("%b", data)
	var i = 0
	for i < len(str) {
		if str[i] != '1' {
			break
		}
		i++
	}
	return i
}

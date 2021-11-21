package column

import "strings"

func Align(s string, repet int) string {
	width := repet*4 - 1
	length := len(s)
	if length < width {
		return s + strings.Repeat("\t", (width-length)/4+1)
	} else {
		s = s[:width]
		return s + "\t"
	}
}

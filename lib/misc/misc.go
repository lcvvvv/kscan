package misc

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
)

func IsDuplicate[T any](slice []T, val T) bool {
	for _, item := range slice {
		if fmt.Sprint(item) == fmt.Sprint(val) {
			return true
		}
	}
	return false
}

func RemoveDuplicateElement[T any](slice []T, elems ...T) []T {
	slice = append(slice, elems...)
	set := make(map[string]struct{}, len(slice))
	j := 0
	for _, v := range slice {
		_, ok := set[fmt.Sprint(v)]
		if ok {
			continue
		}
		set[fmt.Sprint(v)] = struct{}{}
		slice[j] = v
		j++
	}
	return slice[:j]
}

func FixLine(line string) string {
	line = strings.ReplaceAll(line, "\t", "")
	line = strings.ReplaceAll(line, "\r", "")
	line = strings.ReplaceAll(line, "\n", "")
	line = strings.ReplaceAll(line, "\xc2\xa0", "")
	line = strings.ReplaceAll(line, " ", "")
	return line
}

func Xrange(args ...int) []int {
	var start, stop int
	var step = 1
	var r []int
	switch len(args) {
	case 1:
		stop = args[0]
		start = 0
	case 2:
		start, stop = args[0], args[1]
	case 3:
		start, stop, step = args[0], args[1], args[2]
	default:
		return nil
	}
	if start > stop {
		return nil
	}
	if step < 0 {
		return nil
	}

	for i := start; i <= stop; i += step {
		r = append(r, i)
	}
	return r
}

func MustLength(s string, i int) string {
	if len(s) > i {
		return s[:i]
	}
	return s
}

func Percent(int1 int, int2 int) string {
	float1 := float64(int1)
	float2 := float64(int2)
	f := 1 - float1/float2
	f = f * 100
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func StrRandomCut(s string, length int) string {
	sRune := []rune(s)
	if len(sRune) > length {
		i := rand.Intn(len(sRune) - length)
		return string(sRune[i : i+length])
	} else {
		return s
	}
}

func Base64Encode(keyword string) string {
	input := []byte(keyword)
	encodeString := base64.StdEncoding.EncodeToString(input)
	return encodeString
}

func Base64Decode(encodeString string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	return string(decodeBytes), err
}

func CloneStrMap(strMap map[string]string) map[string]string {
	newStrMap := make(map[string]string)
	for k, v := range strMap {
		newStrMap[k] = v
	}
	return newStrMap
}

func CloneIntMap(intMap map[int]string) map[int]string {
	newIntMap := make(map[int]string)
	for k, v := range intMap {
		newIntMap[k] = v
	}
	return newIntMap
}

func RandomString(i ...int) string {
	var length int
	var str string
	if len(i) != 1 {
		length = 32
	} else {
		length = i[0]
	}
	Char := "01234567890abcdef"
	for range Xrange(length) {
		j := rand.Intn(len(Char) - 1)
		str += Char[j : j+1]
	}
	return str
}

func Intersection(a []string, b []string) (inter []string) {
	for _, s1 := range a {
		for _, s2 := range b {
			if s1 == s2 {
				inter = append(inter, s1)
			}
		}
	}
	return inter
}

func FixMap(m map[string]string) map[string]string {
	var arr []string
	rm := make(map[string]string)
	for key, value := range m {
		if value == "" {
			continue
		}
		if IsDuplicate(arr, value) {
			if key != "Username" && key != "Password" {
				continue
			}
		}
		arr = append(arr, value)
		rm[key] = value
	}
	return rm
}

func CloneMap(m map[string]string) map[string]string {
	var nm = make(map[string]string)
	for key, value := range m {
		nm[key] = value
	}
	return nm
}

func AutoWidth(s string, length int) int {
	length1 := len(s)
	length2 := len([]rune(s))

	if length1 == length2 {
		return length
	}
	return length - (length1-length2)/2
}

func ToMap(param interface{}) map[string]string {
	t := reflect.TypeOf(param)
	v := reflect.ValueOf(param)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	m := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		// 通过interface方法来获取key所对应的值
		var cell string
		switch s := v.Field(i).Interface().(type) {
		case string:
			cell = s
		case []string:
			cell = strings.Join(s, "; ")
		case int:
			cell = strconv.Itoa(s)
		case Stringer:
			cell = s.String()
		default:
			continue
		}
		m[t.Field(i).Name] = cell
	}
	return m
}

type Stringer interface {
	String() string
}

func CopySlice[T any](slice []T) []T {
	v := make([]T, len(slice))
	copy(v, slice)
	return v
}

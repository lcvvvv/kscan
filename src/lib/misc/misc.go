package misc

import (
	"bufio"
	"encoding/json"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"strconv"
	"strings"
)

//func TypeOf(v interface{}) string {
//	return fmt.Sprintf("%T", v)
//}

func StrArr2IntArr(strArr []string) ([]int, error) {
	var intArr []int
	for _, value := range strArr {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		intArr = append(intArr, intValue)
	}
	return intArr, nil
}

func StrConcat(s1 string, v ...string) string {
	s2 := strings.Join(v, "")
	return strings.Join([]string{s1, s2}, "")
}

func Str2Int(str string) int {
	intValue, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return intValue
}

//func IntArr2StrArr(intArr []int) []string {
//	var strArr []string
//	for _, value := range intArr {
//		strValue := strconv.Itoa(value)
//		strArr = append(strArr, strValue)
//	}
//	return strArr
//}

func Int2Str(Int int) string {
	return strconv.Itoa(Int)
}

func IsInStrArr(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func IsInIntArr(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func ReadLine(fileName string, handler func(string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		line = FixLine(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func FixLine(line string) string {
	line = strings.Replace(line, "\r", "", -1)
	line = strings.Replace(line, " ", "", -1)
	line = strings.Replace(line, "\t", "", -1)
	line = strings.Replace(line, "\r", "", -1)
	line = strings.Replace(line, "\n", "", -1)
	line = strings.Replace(line, "\xc2", "", -1)
	line = strings.Replace(line, "\xa0", "", -1)
	return line
}

func FillLine(line string) string {
	var length int
	width, _, _ := terminal.GetSize(0)
	width = width - 3
	if len(line) < width {
		length = width - len(line)
	} else {
		length = 0
	}
	return StrConcat(line, strings.Repeat(" ", length))
}

func UniStrAppend(slice []string, elems ...string) []string {
	for _, elem := range elems {
		if IsInStrArr(slice, elem) {
			continue
		} else {
			slice = append(slice, elem)
		}
	}
	return slice
}

func Interface2Str(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

func FileIsExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func SafeOpen(path string) *os.File {
	if FileIsExist(path) {
		f, _ := os.Open(path)
		return f
	} else {
		f, _ := os.Create(path)
		return f
	}
}

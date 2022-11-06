package sflag

import (
	"flag"
	"fmt"
	"io"
	"kscan/lib/misc"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	autoStringSlice      []string
	stringSpliceParamMap = make(map[*string]*[]string)
	intSpliceParamMap    = make(map[*string]*[]int)
	intRangeRegx         = regexp.MustCompile("^([0-9]+)-([0-9]+)$")
	proxyStrRegx         = regexp.MustCompile("^(http|https|socks5|socks4)://[0-9.]+:[0-9]+$")
	NetlocRegx           = regexp.MustCompile("^([.A-Za-z0-9-]+):\\d+$")
)

func BoolVar(p *bool, name string, value bool) {
	flag.BoolVar(p, name, value, "")
}

func StringVar(p *string, name string, value string) {
	flag.StringVar(p, name, value, "")
}

func IntVar(p *int, name string, value int) {
	flag.IntVar(p, name, value, "")
}

func StringSpliceVar(p *[]string, name string) {
	for keyPtr, ptr := range stringSpliceParamMap {
		if ptr == p {
			flag.StringVar(keyPtr, name, "", "")
			return
		}
	}
	var s string
	flag.StringVar(&s, name, "", "")
	stringSpliceParamMap[&s] = p
}

func IntSpliceVar(p *[]int, name string) {
	for keyPtr, ptr := range intSpliceParamMap {
		if ptr == p {
			flag.StringVar(keyPtr, name, "", "")
			return
		}
	}
	var s string
	flag.StringVar(&s, name, "", "")
	intSpliceParamMap[&s] = p
}

func AutoVarString(p *string, name string, value string) {
	flag.StringVar(p, name, value, "")
	autoStringSlice = append(autoStringSlice, "--"+name)
}

func SetUsage(s string) {
	flag.Usage = func() {
		fmt.Print(s)
	}
}

func Parse() {
	fixArgs()
	flag.Parse()
	for strPtr, splicePtr := range stringSpliceParamMap {
		*splicePtr = parseStringSpliceExpress(*strPtr)
	}
	for intPtr, splicePtr := range intSpliceParamMap {
		*splicePtr = parseIntSpliceExpress(*intPtr)
	}
}

func parseStringSpliceExpress(expr string) (splice []string) {
	if expr == "" {
		return splice
	}
	//判断对象是否为多个
	if s := strings.ReplaceAll(expr, "\\,", "[DouHao]"); strings.Count(s, ",") > 0 {
		for _, str := range strings.Split(s, ",") {
			splice = append(splice, strings.ReplaceAll(str, "[DouHao]", ","))
		}
		return splice
	}
	//判断target字符串是否为文件
	if regexp.MustCompile("^file:.+$").MatchString(expr) {
		expr = strings.Replace(expr, "file:", "", 1)
	}
	if _, err := os.Lstat(expr); err == nil {
		fs, err := os.Open(expr)
		if err != nil {
			panic(err)
		}
		buf, err := io.ReadAll(fs)
		if err != nil {
			panic(err)
		}
		splice = strings.Split(string(buf), "\n")
		for i, s := range splice {
			splice[i] = strings.TrimSpace(s)
		}
		return splice
	}
	return append(splice, expr)
}

func parseIntSpliceExpress(expr string) (splice []int) {
	for _, sv := range parseStringSpliceExpress(expr) {
		iv, err := strconv.Atoi(sv)
		if err == nil {
			splice = append(splice, iv)
			continue
		}
		if intRangeRegx.MatchString(sv) == true {
			regRes := intRangeRegx.FindStringSubmatch(sv)
			start, _ := strconv.Atoi(regRes[1])
			end, _ := strconv.Atoi(regRes[2])
			if end > start {
				for j := start; j <= end; j++ {
					splice = append(splice, j)
				}
				continue
			}
		}
		panic(err)
	}
	return splice
}

func NetlocVerification(s string) bool {
	return NetlocRegx.MatchString(s)
}

func ProxyStrVerification(s string) bool {
	return proxyStrRegx.MatchString(s)
}

func fixArgs() {
	var newArgs []string
	for index, value := range os.Args {
		newArgs = append(newArgs, value)
		if misc.IsDuplicate(autoStringSlice, value) {
			if index+2 > len(os.Args) {
				newArgs = append(newArgs, "")
				break
			}
			if os.Args[index+1][:2] == "--" {
				newArgs = append(newArgs, "")
			}
		}
	}
	os.Args = newArgs
}

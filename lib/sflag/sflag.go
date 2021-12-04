package sflag

import (
	"flag"
	"fmt"
	"kscan/lib/misc"
	"os"
	"regexp"
)

var (
	autoStringSlice []string
	multipleIntRegx = regexp.MustCompile("^((?:[0-9]+)(?:-[0-9]+)?)(?:,(?:[0-9]+)(?:-[0-9]+)?)*$")
	multipleStrRegx = regexp.MustCompile("^([.A-Za-z0-9/]+)(,[.A-Za-z0-9/])*$")
	proxyStrRegx    = regexp.MustCompile("^(http|https|socks5|socks4)://[0-9.]+:[0-9]+$")
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
}

func MultipleIntVerification(s string) bool {
	return multipleIntRegx.MatchString(s)
}

func MultipleStrVerification(s string) bool {
	return multipleStrRegx.MatchString(s)

}

func ProxyStrVerification(s string) bool {
	return proxyStrRegx.MatchString(s)

}

func fixArgs() {
	var newArgs []string
	for index, value := range os.Args {
		newArgs = append(newArgs, value)
		if misc.IsInStrArr(autoStringSlice, value) {
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

package fofa

import (
	"fmt"
	"kscan/lib/misc"
	"reflect"
	"regexp"
	"strings"
)

type Result struct {
	Host, Title, Ip, Domain, Port, Country string
	Province, City, Country_name, Protocol string
	Server, Banner, Isp, As_organization   string
	Header, Cert                           string
}

func (r *Result) Fix() {
	if r.Protocol != "" {
		r.Host = fmt.Sprintf("%s://%s:%s", r.Protocol, r.Ip, r.Port)
	}
	if regexp.MustCompile("http([s]?)://.*").MatchString(r.Host) == false && r.Protocol == "" {
		r.Host = "http://" + r.Host
	}
	if r.Banner != "" {
		r.Banner = misc.FixLine(r.Banner)
		r.Banner = misc.StrRandomCut(r.Banner, 20)
	}
	if r.Title == "" && r.Protocol != "" {
		r.Title = strings.ToUpper(r.Protocol)
	}

	r.Title = misc.FixLine(r.Title)

}

func (r Result) Map() map[string]string {
	t := reflect.TypeOf(r)
	v := reflect.ValueOf(r)
	m := make(map[string]string)
	for k := 0; k < t.NumField(); k++ {
		key := t.Field(k).Name
		value := v.Field(k).String()
		m[key] = value
	}
	return m
}

//
//func (r *Result) Output() {
//	slog.Info("搜索关键字为：", r.query)
//	//for _, result := range r.result {
//	//	var row string
//	//	var diff int
//	//	for _, col := range r.field {
//	//		col = misc.FixLine(col)
//	//		col = misc.FilterPrintStr(col)
//	//
//	//		cell := result[col]
//	//		if col == "host" {
//	//			if regexp.MustCompile("http(s?)://.*").MatchString(cell) == false {
//	//				cell = "http://" + cell
//	//			}
//	//		}
//	//		colRuneBuf := []rune(cell)
//	//		length := len(cell)
//	//		width := lengthMap[col]
//	//
//	//		if length >= width && col != "host" {
//	//			cell = string(colRuneBuf[:width-5]) + "..."
//	//		}
//	//		if length+diff >= width {
//	//			diff = length + diff - width
//	//		}
//	//		row = row + fmt.Sprintf("%-"+strconv.Itoa(width)+"v ", cell)
//	//	}
//	//	slog.Data(row)
//	//}
//	t := table.Table(r.result)
//	fmt.Println(t)
//	slog.Info("搜索返回总数量为：", r.size, ",本次返回数量为：", len(r.result))
//}

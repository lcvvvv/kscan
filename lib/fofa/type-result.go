package fofa

import (
	"fmt"
	"kscan/lib/misc"
	"kscan/lib/slog"
	"strconv"
)

var length = map[string]int{
	"host":    28,
	"title":   20,
	"ip":      17,
	"domain":  18,
	"port":    6,
	"country": 3,

	"province":        8,
	"city":            16,
	"country_name":    8,
	"protocol":        8,
	"server":          14,
	"banner":          30,
	"isp":             8,
	"as_organization": 8,

	"header": 30,
	"cert":   30,
}

type Result struct {
	result []map[string]string
	field  []string
	query  string
	size   int
}

func (r *Result) Output() {
	slog.Data("搜索关键字为：", r.query)
	for _, result := range r.result {
		var row string
		for _, col := range r.field {
			col = misc.FixLine(col)
			col = misc.FilterPrintStr(col)

			cell := result[col]
			colRuneBuf := []rune(cell)
			Length := len(cell)
			width := length[col]

			if Length > width {
				cell = string(colRuneBuf[:width-3]) + "..."
			}

			row = row + fmt.Sprintf("%-"+strconv.Itoa(width)+"v", cell)
		}
		slog.Data(row)
	}
	slog.Info("搜索返回总数量为：", r.size, ",本次返回数量为：", len(r.result))
}

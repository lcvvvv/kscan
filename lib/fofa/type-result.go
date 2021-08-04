package fofa

import (
	"kscan/lib/misc"
	"kscan/lib/slog"
)

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
			row = row + result[col] + "\t"
		}
		slog.Data(row)
	}
	slog.Data("搜索返回总数量为：", r.size, ",本次返回数量为：", len(r.result))
}

// Copyright 2017 modood. All rights reserved.
// license that can be found in the LICENSE file.

// Package table produces a string that represents slice of structs data in a text table
package table

import (
	"errors"
	"fmt"
	"reflect"
)

type bd struct {
	H  rune // BOX DRAWINGS HORIZONTAL
	V  rune // BOX DRAWINGS VERTICAL
	VH rune // BOX DRAWINGS VERTICAL AND HORIZONTAL
	HU rune // BOX DRAWINGS HORIZONTAL AND UP
	HD rune // BOX DRAWINGS HORIZONTAL AND DOWN
	VL rune // BOX DRAWINGS VERTICAL AND LEFT
	VR rune // BOX DRAWINGS VERTICAL AND RIGHT
	DL rune // BOX DRAWINGS DOWN AND LEFT
	DR rune // BOX DRAWINGS DOWN AND RIGHT
	UL rune // BOX DRAWINGS UP AND LEFT
	UR rune // BOX DRAWINGS UP AND RIGHT
}

var (
	DrawingBox   = true
	PrintColumns []string
)

var m = map[string]bd{
	"ascii":       {'-', '|', '+', '+', '+', '+', '+', '+', '+', '+', '+'},
	"box-drawing": {'─', '│', '┼', '┴', '┬', '┤', '├', '┐', '┌', '┘', '└'},
}

// Output formats slice of structs data and writes to standard output.(Using box drawing characters)
func Output(slice interface{}) {
	fmt.Println(Table(slice))
}

// OutputA formats slice of structs data and writes to standard output.(Using standard ascii characters)
func OutputA(slice interface{}) {
	fmt.Println(AsciiTable(slice))
}

// Table formats slice of structs data and returns the resulting string.(Using box drawing characters)
func Table(slice interface{}) string {
	coln, colw, rows, err := parse(slice)
	if err != nil {
		return err.Error()
	}
	table := table(coln, colw, rows, m["box-drawing"])
	return table
}

// AsciiTable formats slice of structs data and returns the resulting string.(Using standard ascii characters)
func AsciiTable(slice interface{}) string {
	coln, colw, rows, err := parse(slice)
	if err != nil {
		return err.Error()
	}
	table := table(coln, colw, rows, m["ascii"])
	return table
}

func parse(slice interface{}) (
	coln []string, // name of columns
	colw []int, // width of columns
	rows [][]string, // rows of content
	err error,
) {

	s, err := sliceconv(slice)
	if err != nil {
		return
	}
	for i, u := range s {
		v := reflect.ValueOf(u)
		t := reflect.TypeOf(u)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
			t = t.Elem()
		}
		if v.Kind() != reflect.Struct {
			err = errors.New("warning: table: items of slice should be on struct value")
			return
		}
		var row []string

		m := 0 // count of unexported field
		for n := 0; n < v.NumField(); n++ {
			if t.Field(n).PkgPath != "" {
				m++
				continue
			}
			cn := t.Field(n).Name
			ct := t.Field(n).Tag.Get("table")
			if ct == "" {
				ct = cn
			}

			if len(PrintColumns) > 0 {
				b := false
				for _, name := range PrintColumns {
					if name == ct {
						b = true
					}
				}
				if b == false {
					m++
					continue
				}
			}

			cv := fmt.Sprintf("%+v", v.FieldByName(cn).Interface())

			if i == 0 {
				coln = append(coln, ct)
				colw = append(colw, len(ct))
			}
			if colw[n-m] < getLength(cv) {
				colw[n-m] = getLength(cv)
			}

			row = append(row, cv)
		}
		rows = append(rows, row)
	}
	return coln, colw, rows, nil
}

func table(coln []string, colw []int, rows [][]string, b bd) (table string) {
	if len(rows) == 0 {
		return ""
	}

	head := [][]rune{{b.DR}, {b.V}, {b.VR}}
	bttm := []rune{b.UR}
	for i, v := range colw {
		head[0] = append(head[0], []rune(repeat(v+2, b.H)+string(b.HD))...)
		head[1] = append(head[1], []rune(" "+coln[i]+repeat(v-StringLength([]rune(coln[i]))+1, ' ')+string(b.V))...)
		head[2] = append(head[2], []rune(repeat(v+2, b.H)+string(b.VH))...)
		bttm = append(bttm, []rune(repeat(v+2, b.H)+string(b.HU))...)
	}
	head[0][len(head[0])-1] = b.DL
	head[2][len(head[2])-1] = b.VL
	bttm[len(bttm)-1] = b.UL

	var body [][]rune
	for _, r := range rows {
		row := []rune{b.V}
		for i, v := range colw {
			// handle non-ascii character
			l := StringLength([]rune(r[i]))

			row = append(row, []rune(" "+r[i]+repeat(v-l+1, ' ')+string(b.V))...)
		}
		body = append(body, row)
	}

	for _, v := range head {
		table += string(v) + "\n"
	}
	for _, v := range body {
		table += string(v) + "\n"
	}
	table += string(bttm)
	return table
}

func sliceconv(slice interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("warning: sliceconv: param \"slice\" should be on slice value")
	}

	l := v.Len()
	r := make([]interface{}, l)
	for i := 0; i < l; i++ {
		r[i] = v.Index(i).Interface()
	}
	return r, nil
}

func repeat(time int, char rune) string {
	var s = make([]rune, time)
	for i := range s {
		s[i] = char
	}
	return string(s)
}

// StringLength string display length
func StringLength(r []rune) int {
	// CJK(Chinese, Japanese, Korean)
	type cjk struct {
		from rune
		to   rune
	}

	// References:
	// -   [Unicode Table](http://www.tamasoft.co.jp/en/general-info/unicode.html)
	// -   [汉字 Unicode 编码范围](http://www.qqxiuzi.cn/zh/hanzi-unicode-bianma.php)

	var a = []cjk{
		{0x2E80, 0x9FD0},   // Chinese, Hiragana, Katakana, ...
		{0xAC00, 0xD7A3},   // Hangul
		{0xF900, 0xFACE},   // Kanji
		{0xFE00, 0xFE6C},   // Fullwidth
		{0xFF00, 0xFF60},   // Fullwidth again
		{0x20000, 0x2FA1D}, // Extension
		// More? PRs are aways welcome here.
	}
	length := len(r)
l:
	for _, v := range r {
		for _, c := range a {
			if v >= c.from && v <= c.to {
				length++
				continue l
			}
		}
	}
	return length
}

func SetDrawingBox(b bool) {
	DrawingBox = b
}

func SetPrintColumns(s []string) {
	PrintColumns = s
}

func getLength(s string) int {
	length1 := len(s)
	runeBuf := []rune(s)
	length2 := len(runeBuf)

	return (length1-length2)/2 + length2
}

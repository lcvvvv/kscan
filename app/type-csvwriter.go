package app

import "encoding/csv"

type CSVWriter struct {
	f     *csv.Writer
	title []string
}

func (cw *CSVWriter) inTitle(title string) bool {
	for _, value := range cw.title {
		if value == title {
			return true
		}
	}
	return false
}

func (cw *CSVWriter) Push(m map[string]string) {
	var cells []string
	for _, key := range cw.title {
		if value, ok := m[key]; ok {
			cells = append(cells, value)
			delete(m, key)
		} else {
			cells = append(cells, "")
		}
	}
	for key, value := range m {
		cells = append(cells, value)
		cw.title = append(cw.title, key)
	}
	cw.f.Write(cells)
	cw.f.Flush()
}

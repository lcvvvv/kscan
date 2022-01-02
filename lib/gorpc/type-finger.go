package gorpc

import "strings"

type Finger struct {
	data []string
}

func New() *Finger {
	return &Finger{
		data: []string{},
	}
}

func (f *Finger) Append(s string) {
	f.data = append(f.data, s)
}

func (f *Finger) Value() []string {
	return f.data
}

func (f *Finger) ValueString() string {
	return strings.Join(f.data, ",")
}

func (f *Finger) Len() int {
	return len(f.data)
}

package app

import (
	"encoding/json"
	"kscan/core/slog"
	"os"
	"sync"
)

type JSONWriter struct {
	f     *os.File
	mutex *sync.Mutex
}

func (jw *JSONWriter) Push(m map[string]string) {
	jw.mutex.Lock()
	defer jw.mutex.Unlock()
	stat, err := jw.f.Stat()
	if err != nil {
		slog.Println(slog.ERROR, err)
	}
	jsonBuf, _ := json.MarshalIndent(m, "\t", "\t")
	jsonBuf = append(jsonBuf, []byte("\n]\n")...)
	if stat.Size() == 2 {
		jw.f.Seek(stat.Size()-1, 0)
		jsonBuf = append([]byte("\n\t"), jsonBuf...)
		jw.f.Write(jsonBuf)
	} else {
		jw.f.Seek(stat.Size()-4, 0)
		jsonBuf = append([]byte("},\n\t"), jsonBuf...)
		jw.f.Write(jsonBuf)
	}
}

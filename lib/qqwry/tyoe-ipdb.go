//https://github.com/zu1k/nali
package qqwry

import (
	"fmt"
	"os"
)

// FileData: info of database file
type FileData struct {
	Data     []byte
	FilePath string
	FileBase *os.File
}

// IPDB common ip database
type IPDB struct {
	Data   *FileData
	Offset uint32
	IPNum  uint32
}

// setOffset 设置偏移量
func (db *IPDB) SetOffset(offset uint32) {
	db.Offset = offset
}

// readString 获取字符串
func (db *IPDB) ReadString(offset uint32) []byte {
	db.SetOffset(offset)
	data := make([]byte, 0, 30)
	buf := make([]byte, 1)
	for {
		buf = db.ReadData(1)
		if buf[0] == 0 {
			break
		}
		data = append(data, buf[0])
	}
	return data
}

// readData 从文件中读取数据
func (db *IPDB) ReadData(length uint32, offset ...uint32) (rs []byte) {
	if len(offset) > 0 {
		db.SetOffset(offset[0])
	}

	end := db.Offset + length
	dataNum := uint32(len(db.Data.Data))
	if db.Offset > dataNum {
		return nil
	}

	if end > dataNum {
		end = dataNum
	}
	rs = db.Data.Data[db.Offset:end]
	db.Offset = end
	return
}

// readMode 获取偏移值类型
func (db *IPDB) ReadMode(offset uint32) byte {
	mode := db.ReadData(1, offset)
	return mode[0]
}

// ReadUInt24
func (db *IPDB) ReadUInt24() uint32 {
	buf := db.ReadData(3)
	return ByteToUInt32(buf)
}

// readArea 读取区域
func (db *IPDB) ReadArea(offset uint32) []byte {
	mode := db.ReadMode(offset)
	if mode == RedirectMode1 || mode == RedirectMode2 {
		areaOffset := db.ReadUInt24()
		if areaOffset == 0 {
			return []byte("")
		}
		return db.ReadString(areaOffset)
	}
	return db.ReadString(offset)
}

func GetMiddleOffset(start uint32, end uint32, indexLen uint32) uint32 {
	records := ((end - start) / indexLen) >> 1
	return start + records*indexLen
}

type Result struct {
	Country string
	Area    string
}

func (r Result) String() string {
	return fmt.Sprintf("%s %s", r.Country, r.Area)
}

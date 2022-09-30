//https://github.com/zu1k/nali
package qqwry

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

type QQwry struct {
	IPDB
}

var logger = Logger(log.New(os.Stderr, "[qqwry]", log.Ldate|log.Ltime))

type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}

const (
	// RedirectMode1 [IP][0x01][国家和地区信息的绝对偏移地址]
	RedirectMode1 = 0x01
	// RedirectMode2 [IP][0x02][信息的绝对偏移][...] or [IP][国家][...]
	RedirectMode2 = 0x02
)

// NewQQwry new database from path
func NewQQwryFS(fs *os.File) (qqwry *QQwry, err error) {
	var fileData []byte
	var fileInfo FileData

	fileInfo.FileBase = fs
	defer fileInfo.FileBase.Close()

	fileData, err = ioutil.ReadAll(fileInfo.FileBase)
	if err != nil {
		return nil, err
	}

	fileInfo.Data = fileData

	buf := fileInfo.Data[0:8]
	start := binary.LittleEndian.Uint32(buf[:4])
	end := binary.LittleEndian.Uint32(buf[4:])

	return &QQwry{
		IPDB: IPDB{
			Data:  &fileInfo,
			IPNum: (end-start)/7 + 1,
		},
	}, nil
}

func (db QQwry) Find(query string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("query keyword should be IPv4")
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, errors.New("query keyword should be IPv4")
	}
	ip4uint := binary.BigEndian.Uint32(ip4)

	offset := db.searchIndex(ip4uint)
	if offset <= 0 {
		return nil, errors.New("query keyword not valid")
	}

	var gbkCountry []byte
	var gbkArea []byte

	mode := db.ReadMode(offset + 4)
	switch mode {
	case RedirectMode1: // [IP][0x01][国家和地区信息的绝对偏移地址]
		countryOffset := db.ReadUInt24()
		mode = db.ReadMode(countryOffset)
		if mode == RedirectMode2 {
			c := db.ReadUInt24()
			gbkCountry = db.ReadString(c)
			countryOffset += 4
		} else {
			gbkCountry = db.ReadString(countryOffset)
			countryOffset += uint32(len(gbkCountry) + 1)
		}
		gbkArea = db.ReadArea(countryOffset)
	case RedirectMode2:
		countryOffset := db.ReadUInt24()
		gbkCountry = db.ReadString(countryOffset)
		gbkArea = db.ReadArea(offset + 8)
	default:
		gbkCountry = db.ReadString(offset + 4)
		gbkArea = db.ReadArea(offset + uint32(5+len(gbkCountry)))
	}

	enc := simplifiedchinese.GBK.NewDecoder()
	country, _ := enc.String(string(gbkCountry))
	area, _ := enc.String(string(gbkArea))
	result = Result{
		Country: strings.ReplaceAll(country, " CZ88.NET", ""),
		Area:    strings.ReplaceAll(area, " CZ88.NET", ""),
	}
	return result, nil
}

// searchIndex 查找索引位置
func (db *QQwry) searchIndex(ip uint32) uint32 {
	header := db.ReadData(8, 0)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	buf := make([]byte, 7)
	mid := uint32(0)
	ipUint := uint32(0)

	for {
		mid = GetMiddleOffset(start, end, 7)
		buf = db.ReadData(7, mid)
		ipUint = binary.LittleEndian.Uint32(buf[:4])

		if end-start == 7 {
			offset := ByteToUInt32(buf[4:])
			buf = db.ReadData(7)
			if ip < binary.LittleEndian.Uint32(buf[:4]) {
				return offset
			}
			return 0
		}

		if ipUint > ip {
			end = mid
		} else if ipUint < ip {
			start = mid
		} else if ipUint == ip {
			return ByteToUInt32(buf[4:])
		}
	}
}

func ByteToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}

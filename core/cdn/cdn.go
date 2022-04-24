package cdn

import (
	"kscan/core/slog"
	"kscan/lib/qqwry"
	"os"
	"path/filepath"
	"strings"
)

var database *qqwry.QQwry

var filename = "qqwry.dat"
var path = getRealPath()

func Init() {
	d, err := qqwry.NewQQwry(GetPath())
	if err != nil {
		slog.Println(slog.WARN, "qqwry init err:", err)
		return
	}
	database = d
}

func GetPath() string {
	return path + "/" + filename
}

func DownloadQQWry() error {
	return qqwry.Download(GetPath())
}

func FindCDN(query string) (bool, string, error) {
	result, err := Find(query)
	if strings.Contains(result, "CDN") {
		return true, result, err
	}
	return false, "", err
}

func Find(query string) (string, error) {
	result, err := database.Find(query)
	if err != nil {
		return "", err
	}
	return result.String(), err
}

func getRealPath() string {
	dir, _ := os.Executable()
	path := filepath.Dir(dir)
	return path
}

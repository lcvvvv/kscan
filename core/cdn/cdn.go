package cdn

import (
	"kscan/core/slog"
	"kscan/lib/qqwry"
	"strings"
)

var database *qqwry.QQwry

const filename = "./qqwry.dat"

func Init() {
	d, err := qqwry.NewQQwry(filename)
	if err != nil {
		slog.Println(slog.WARN, "qqwry init err:", err)
		return
	}
	database = d
}

func DownloadQQWry() error {
	return qqwry.Download(filename)
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

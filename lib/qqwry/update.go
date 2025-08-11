// https://github.com/zu1k/nali
package qqwry

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const mirror = "https://qqwry-mirror.cdn.skk.moe/qqwry.dat"

func Download(filePath string) error {
	data, err := Get(mirror)
	if err != nil {
		return err
	}
	return SaveFile(filePath, data)
}

func Get(url string) ([]byte, error) {
	var UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36"
	Client := http.DefaultClient
	Client.Timeout = time.Second * 30
	Client.Transport = &http.Transport{
		TLSHandshakeTimeout:   time.Second * 5,
		IdleConnTimeout:       time.Second * 20,
		ResponseHeaderTimeout: time.Second * 20,
		ExpectContinueTimeout: time.Second * 20,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("User-Agent", UserAgent)
	resp, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("http response is nil")
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("http response status code is not 200")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func SaveFile(path string, data []byte) (err error) {
	// Remove file if exist
	_, err = os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			logger.Println("旧文件删除失败", err.Error())
		}
	}

	// save file
	return ioutil.WriteFile(path, data, 0644)
}

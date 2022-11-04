package telnet

import (
	"kscan/lib/gotelnet"
)

func Check(addr, username, password string, port, serverType int) error {
	client := gotelnet.New(addr, port)
	err := client.Connect()
	if err != nil {
		return err
	}
	defer client.Close()
	client.UserName = username
	client.Password = password
	client.ServerType = serverType
	err = client.Login()
	if err != nil {
		return err
	}
	return nil
}

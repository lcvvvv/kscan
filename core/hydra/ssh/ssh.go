package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

func Check(Host, Username, Password string, Port int) error {
	var Auth = []ssh.AuthMethod{ssh.Password(Password)}
	config := &ssh.ClientConfig{
		User:    Username,
		Auth:    Auth,
		Timeout: 3 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", Host, Port), config)
	if err != nil {
		return err
	}
	defer client.Close()
	return nil
}

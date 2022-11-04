package smb

import (
	"context"
	"errors"
	"github.com/stacktitan/smb/smb"
	"time"
)

var (
	LoginFailedError  = errors.New("login failed")
	LoginTimeoutError = errors.New("login timeout")
)

func Check(Host, Username, Domain, Password string, Port int) error {
	status := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	options := smb.Options{
		Host:        Host,
		Port:        Port,
		User:        Username,
		Password:    Password,
		Domain:      Domain,
		Workstation: "",
	}
	//开始进行SMB连接
	go func() {
		session, err := smb.NewSession(options, false)
		if err != nil {
			status <- err
			return
		}
		defer session.Close()
		if session.IsAuthenticated == false {
			status <- LoginFailedError
			return
		}
		status <- nil
	}()

	select {
	case <-ctx.Done():
		return LoginTimeoutError
	case err := <-status:
		return err
	}

}

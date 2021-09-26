package telnet

import (
	"fmt"
	"net"
	"time"
)

const (
	TIME_DELAY_AFTER_WRITE = 300 * time.Millisecond
)

type Client struct {
	IPAddr   string
	Port     int
	UserName string
	Password string
	conn     net.Conn
}

func New(addr string, port int) *Client {
	return &Client{IPAddr: addr, Port: port}
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.Netloc(), 5*time.Second)
	if err != nil {
		return err
	}
	c.conn = conn
	c.initConnect()
	return nil
}

func (c *Client) Read() ([]byte, error) {
	var buf [2048]byte
	var n int
	_ = c.conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	//读取初始化数据
	n, err := c.conn.Read(buf[0:])
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func (c *Client) Write(buf []byte) error {
	_ = c.conn.SetWriteDeadline(time.Now().Add(time.Second * 3))
	_, err := c.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) analyzeNVT(buf []byte) ([]byte, string) {
	var indexArr []int
	var nvt []byte
	var data string
	fmt.Println(buf)
	for index, b := range buf {
		if b == 255 {
			indexArr = append(indexArr, index)
		}
	}
	for i, index := range indexArr {
		if len(buf) > index+2 {
			nvt = append(nvt, buf[index], buf[index+1], buf[index+2])
		}
		if i == 0 {
			continue
		}
		fmt.Println(indexArr[i-1])
		fmt.Println(indexArr[i])
		startIndex := indexArr[i-1] + 2
		stopIndex := indexArr[i]
		data += string(buf[startIndex:stopIndex])
	}
	return buf, data
}

func (c *Client) initConnect() {
	var buf []byte
	for {
		buf, _ = c.Read()
		nvt, _ := c.analyzeNVT(buf)
		if len(nvt) == 0 {
			break
		}
		var WriteBuf []byte
		for i := 0; i <= len(nvt)/3; i++ {
			start := i * 3
			command := NewNvt(nvt[start : start+2])
			WriteBuf = append(WriteBuf, command.DO()...)
		}
		_ = c.Write(WriteBuf)
	}

}

func (c *Client) Netloc() string {
	return fmt.Sprintf("%s:%d", c.IPAddr, c.Port)
}

func (c *Client) Close() {
	c.conn.Close()
}

//
//func Check(Host, Username, Password string, Port int) (bool, error) {

//	//空行，比对返回数据
//	n, err = conn.Write([]byte("\n"))
//	if nil != err {
//		return false, err
//	}
//	n, err = conn.Read(buf[0:])
//	time.Sleep(500 * time.Millisecond)
//	if nil != err {
//		return false, err
//	}
//	result = string(buf[:n])
//	fmt.Println(netloc, "|", result)
//	return false, err
//	//_, err = conn.Write([]byte(fmt.Sprintf("auth %s\r\n", Password)))
//	//time.Sleep(time.Millisecond * 500)
//	//if err != nil {
//	//	return false, err
//	//}
//	//reply, err := readResponse(conn)
//	//if err != nil {
//	//	return false, err
//	//}
//	//if strings.Contains(reply, "+OK") == false {
//	//	return false, err
//	//}
//	//return true, err
//}
//
////
////func (this *TelnetClient) telnetProtocolHandshake(conn net.Conn) bool {
////	var buf [4096]byte
////	n, err := conn.Read(buf[0:])
////	if nil != err {
////		return false
////	}
////	fmt.Println(string(buf[0:n]))
////	fmt.Println(buf[0:n])
////
////	buf[1] = 252
////	buf[4] = 252
////	buf[7] = 252
////	buf[10] = 252
////	fmt.Println((buf[0:n]))
////	n, err = conn.Write(buf[0:n])
////	if nil != err {
////		log.Error("pkg: model, func: telnetProtocolHandshake, method: conn.Write, errInfo:", err)
////		return false
////	}
////	time.Sleep(time.Millisecond * TIME_DELAY_AFTER_WRITE)
////
////	n, err = conn.Read(buf[0:])
////	if nil != err {
////		log.Error("pkg: model, func: telnetProtocolHandshake, method: conn.Read, errInfo:", err)
////		return false
////	}
////	fmt.Println(string(buf[0:n]))
////	fmt.Println((buf[0:n]))
////
////	buf[1] = 252
////	buf[4] = 251
////	buf[7] = 252
////	buf[10] = 254
////	buf[13] = 252
////	fmt.Println((buf[0:n]))
////	n, err = conn.Write(buf[0:n])
////	if nil != err {
////		log.Error("pkg: model, func: telnetProtocolHandshake, method: conn.Write, errInfo:", err)
////		return false
////	}
////	time.Sleep(time.Millisecond * TIME_DELAY_AFTER_WRITE)
////
////	n, err = conn.Read(buf[0:])
////	if nil != err {
////		log.Error("pkg: model, func: telnetProtocolHandshake, method: conn.Read, errInfo:", err)
////		return false
////	}
////	fmt.Println(string(buf[0:n]))
////	fmt.Println((buf[0:n]))
////
////	if false == this.IsAuthentication {
////		return true
////	}
////
////	n, err = conn.Write([]byte(this.UserName + "\n"))
////	if nil != err {
////		log.Error("pkg: model, func: telnetProtocolHandshake, method: conn.Write, errInfo:", err)
////		return false
////	}
////	time.Sleep(time.Millisecond * TIME_DELAY_AFTER_WRITE)
////
////	n, err = conn.Read(buf[0:])
////	if nil != err {
////		log.Error("pkg: model, func: telnetProtocolHandshake, method: conn.Read, errInfo:", err)
////		return false
////	}
////	fmt.Println(string(buf[0:n]))
////
////	n, err = conn.Write([]byte(this.Password + "\n"))
////	if nil != err {
////		log.Error("pkg: model, func: telnetProtocolHandshake, method: conn.Write, errInfo:", err)
////		return false
////	}
////	time.Sleep(time.Millisecond * TIME_DELAY_AFTER_WRITE)
////
////	n, err = conn.Read(buf[0:])
////	if nil != err {
////		log.Error("pkg: model, func: telnetProtocolHandshake, method: conn.Read, errInfo:", err)
////		return false
////	}
////	fmt.Println(string(buf[0:n]))
////	return true
////}

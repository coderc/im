package sdk

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/coderc/im/common/tcp"
)

type connect struct {
	sendChan, recvChan chan *Message
	conn               *net.TCPConn
}

func newConnect(ip net.IP, port int) *connect {
	clientConn := &connect{
		sendChan: make(chan *Message),
		recvChan: make(chan *Message),
	}
	addr := &net.TCPAddr{IP: ip, Port: port}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Printf("DialTcp.err=%+v", err)
		return nil
	}
	clientConn.conn = conn
	go func() {
		for {
			data, err := tcp.ReadData(conn)
			if err != nil {
				fmt.Printf("ReadData err: %+v", err)
			}
			msg := &Message{}
			json.Unmarshal(data, msg)
			clientConn.recvChan <- msg
		}
	}()

	return clientConn
}

func (c *connect) send(msg *Message) {
	bytes, _ := json.Marshal(msg)
	dataPkg := tcp.DataPkg{
		Data: bytes,
		Len:  uint32(len(bytes)),
	}
	dataPkgBytes := dataPkg.Marshal()
	c.conn.Write(dataPkgBytes)
}

func (c *connect) recv() <-chan *Message {
	return c.recvChan
}

func (c *connect) close() {
	c.conn.Close()
}

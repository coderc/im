package sdk

type connect struct {
	serverAddr         string
	sendChan, recvChan chan *Message
}

func newConnect(serverAddr string) *connect {
	return &connect{
		serverAddr: serverAddr,
		sendChan:   make(chan *Message),
		recvChan:   make(chan *Message),
	}
}

func (c *connect) send(msg *Message) {
	// 不做任何处理，模拟消息发送操作
	c.recvChan <- msg
}

func (c *connect) recv() <-chan *Message {
	return c.recvChan
}

func (c *connect) close() {
}

package sdk

import "net"

const (
	MsgTypeText = "text"
)

type Chat struct {
	Nick      string
	UserId    string
	SessionId string
	conn      *connect
}

type Message struct {
	Type       string
	Name       string
	FormUserId string
	ToUserId   string
	Content    string
	Session    string
}

func NewChat(ip net.IP, port int, nick, userId, sessionId string) *Chat {
	return &Chat{
		Nick:      nick,
		UserId:    userId,
		SessionId: sessionId,
		conn:      newConnect(ip, port),
	}
}

func (c *Chat) Send(msg *Message) {
	c.conn.send(msg)
}

func (c *Chat) Close() {
	c.conn.close()
}

func (c *Chat) Recv() <-chan *Message {
	return c.conn.recv()
}

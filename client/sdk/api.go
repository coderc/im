package sdk

const (
	MsgTypeText = "text"
)

type Chat struct {
	Nick      string
	UserID    string
	SessionID string
	conn      *connect
}

type Message struct {
	Type       string
	Name       string
	FormUserID string
	ToUserID   string
	Content    string
	Session    string
}

func NewChat(serverAddr, nick, userID, sessionID string) *Chat {
	return &Chat{
		Nick:      nick,
		UserID:    userID,
		SessionID: sessionID,
		conn:      newConnet(serverAddr),
	}
}

func (c *Chat) Send(msg *Message) {
	c.conn.send(msg)
}

func (c *Chat) Recv() <-chan *Message {
	return c.conn.Recv()
}

func (c *Chat) Close() {
	c.conn.close()
}

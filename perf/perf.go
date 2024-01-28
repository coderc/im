package perf

import (
	"fmt"
	"net"

	"github.com/coderc/im/client/sdk"
)

var (
	TcpConnNum int32
)

func Run() {
	for i := 0; i < int(TcpConnNum); i++ {
		sdk.NewChat(net.ParseIP("127.0.0.1"), 8900,
			fmt.Sprintf("nickId-test-%d", i), fmt.Sprintf("userId-test-%d", i), fmt.Sprintf("sessionId-test-%d", i))
	}
}

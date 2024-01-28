package gateway

import (
	"errors"
	"io"
	"log"
	"net"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/coderc/im/common/config"
	"github.com/coderc/im/common/tcp"
)

// Run 网关服务
func Run() {
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{Port: config.GetGatewayServerPort()})
	if err != nil {
		log.Fatalf("StartTCPEpollServer err: %s", err.Error())
		panic(err)
	}
	initWorkPool()
	initEpoll(ln, runProc)

	logger.Infof("---------- im gateway stated ----------")
	select {}
}

func runProc(c *connection, ep *epoller) {
	// 读取一个完整的消息包
	dataBuf, err := tcp.ReadData(c.conn)
	if err != nil {
		// 如果读取conn时连接关闭，则直接端口连接
		if errors.Is(err, io.EOF) {
			ep.remove(c)
		}
		return
	}

	logger.Debug(string(dataBuf))
	err = wPool.Submit(func() {
		// 交给 state server rpc 处理
		bytes := tcp.DataPkg{
			Len:  uint32(len(dataBuf)),
			Data: dataBuf,
		}
		tcp.SendData(c.conn, bytes.Marshal())
	})

	if err != nil {
		logger.Errorf("runProc err: %v", err)
	}
}

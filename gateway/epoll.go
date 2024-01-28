package gateway

import (
	"net"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/coderc/im/common/config"
	"golang.org/x/sys/unix"
)

var (
	ep     *ePool // epoll池
	tcpNum int32  // 当前服务器允许接入的最大tcp的连接数
)

type ePool struct {
	eChan  chan *connection
	tables sync.Map
	eSize  int
	done   chan struct{}

	ln *net.TCPListener
	f  func(c *connection, ep *epoller)
}

func initEpoll(ln *net.TCPListener, f func(c *connection, ep *epoller)) {
	setLimit()
	ep = newEPool(ln, f)
	ep.createAcceptProcess()
	ep.startEPool()
}

func newEPool(ln *net.TCPListener, cb func(c *connection, ep *epoller)) *ePool {
	return &ePool{
		eChan:  make(chan *connection, config.GetGatewayEpollerChanNum()),
		done:   make(chan struct{}),
		eSize:  config.GetGatewayEpollerNum(),
		tables: sync.Map{},
		ln:     ln,
		f:      cb,
	}
}

// 创建一批专门处理accept事件的协程，协程数与当前CPU核心数对应
func (e *ePool) createAcceptProcess() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				conn, err := e.ln.AcceptTCP()
				if !checkTcp() {
					_ = conn.Close()
					continue
				}
				setTcpConfig(conn)
				if err != nil {
					if ne, ok := err.(net.Error); ok && ne.Temporary() {
						logger.Errorf("accept temp err: %v", ne)
						continue
					}
					logger.Errorf("accept err: %v", err)
				}

				c := connection{
					conn: conn,
					fd:   socketFD(conn),
				}
				ep.addTack(&c)
			}
		}()
	}
}

func (e *ePool) startEPool() {
	for i := 0; i < e.eSize; i++ {
		go e.startEProc()
	}
}

func (e *ePool) startEProc() {
	ep, err := newepoller()
	if err != nil {
		panic(err)
	}

	// 监听连接创建事件
	go func() {
		for {
			select {
			case <-e.done: // ePoll关闭
				return
			case conn := <-e.eChan:
				addTcpNum()
				logger.Infof("tcpNum: %d", getTcpNum())
				if err := ep.add(conn); err != nil {
					logger.Warnf("failed to add connection: %v", err)
					conn.Close() // 登陆未成功直接关闭此连接
					continue
				}
				logger.Infof("EpollerPoll new connection[%v] tcpSize: %d", conn.RemoteAddr(), getTcpNum())
			}
		}
	}()

	// 轮询器 轮询等待wait事件触发回调函数
	for {
		select {
		case <-e.done: // ePoll关闭
			return
		default:
			connections, err := ep.wait(200) // 设置轮询的最小间隔，避免忙轮询

			if err != nil && err != syscall.EINTR {
				logger.Warnf("failed to ePoll wait, err: %v", err)
				continue
			}

			for _, conn := range connections {
				if conn == nil {
					break
				}
				e.f(conn, ep)
			}

		}
	}
}

func (e *ePool) addTack(c *connection) {
	e.eChan <- c
}

// epoller 轮询器
type epoller struct {
	fd int
}

func newepoller() (*epoller, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoller{
		fd: fd,
	}, nil
}

func (e *epoller) add(conn *connection) error {
	fd := conn.fd

	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.EPOLLIN | unix.EPOLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}
	ep.tables.Store(fd, conn)
	return nil
}

func (e *epoller) remove(c *connection) error {
	subTcpNum()
	fd := c.fd
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return nil
	}
	ep.tables.Delete(fd)
	return nil
}

func (e *epoller) wait(msec int) ([]*connection, error) {
	events := make([]unix.EpollEvent, config.GetGatewayEpollWaitQueueSize())
	n, err := unix.EpollWait(e.fd, events, msec)
	if err != nil {
		return nil, err
	}
	var connections []*connection
	for i := 0; i < n; i++ {
		if conn, ok := ep.tables.Load(int(events[i].Fd)); ok {
			connections = append(connections, conn.(*connection))
		}
	}

	return connections, nil
}

// socketFD 通过反射拿到私有变量的成员值，从而达到拿到TCP连接fd的目的
func socketFD(conn *net.TCPConn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(*conn).FieldByName("conn"))
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

// setLimit 设置go进程持有文件句柄数限制
func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	logger.Infof("set cur limit: %d", rLimit.Cur)
}

func addTcpNum() {
	atomic.AddInt32(&tcpNum, 1)
}

func getTcpNum() int32 {
	return atomic.LoadInt32(&tcpNum)
}

func subTcpNum() {
	atomic.AddInt32(&tcpNum, -1)
}

func checkTcp() bool {
	num := getTcpNum()
	maxTcpNum := config.GetGatewayMaxTcpNum()
	return num <= maxTcpNum
}

// setTcpConfig 对TCP连接进行设置
func setTcpConfig(c *net.TCPConn) {
	_ = c.SetKeepAlive(true)
}

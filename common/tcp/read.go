package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func ReadData(conn *net.TCPConn) ([]byte, error) {
	var dataLen uint32
	dataLenBuf := make([]byte, 4)
	if err := readFixedData(conn, dataLenBuf); err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(dataLenBuf)
	if err := binary.Read(buffer, binary.BigEndian, &dataLen); err != nil {
		return nil, fmt.Errorf("read headlen err: %v", err)
	}
	if dataLen <= 0 {
		return nil, fmt.Errorf("wrong headlen: %d", dataLen)
	}

	dataBuf := make([]byte, dataLen)
	if err := readFixedData(conn, dataBuf); err != nil {
		return nil, fmt.Errorf("read headlen err: %v", err)
	}

	return dataBuf, nil
}

// readFixedData 读取固定长度的数据到buf中
func readFixedData(conn *net.TCPConn, buf []byte) error {
	_ = (*conn).SetReadDeadline(time.Now().Add(time.Duration(120) * time.Second)) // 设置读取TCP数据包的超时时间
	var pos int = 0
	var totalSize int = len(buf)
	for {
		c, err := (*conn).Read(buf[pos:])
		if err != nil {
			return err
		}
		pos = pos + c
		if pos == totalSize {
			break
		}
	}

	return nil
}

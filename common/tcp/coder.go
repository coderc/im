package tcp

import (
	"bytes"
	"encoding/binary"
)

type DataPkg struct {
	Len  uint32
	Data []byte
}

func (d *DataPkg) Marshal() []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, d.Len)
	return append(bytesBuffer.Bytes(), d.Data...)
}

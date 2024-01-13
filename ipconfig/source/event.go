package source

import (
	"fmt"

	"github.com/coderc/im/common/discovery"
)

var eventChan chan *Event

func EventChan() <-chan *Event {
	return eventChan
}

type EventType string

const (
	AddNodeEvent EventType = "addNode"
	DelNodeEvent EventType = "delNode"
)

type Event struct {
	Type         EventType
	IP           string
	Port         string
	ConnectNum   float64
	MessageBytes float64
}

func NewEvent(ed *discovery.EndPointInfo) *Event {
	if ed == nil || ed.MetaData == nil {
		return nil
	}

	var connNum, msgBytes float64
	if data, ok := ed.MetaData["connect_num"]; ok {
		if connNum, ok = data.(float64); !ok {
			panic("connect_num assert float64 failed")
		}
	}
	if data, ok := ed.MetaData["message_bytes"]; ok {
		if msgBytes, ok = data.(float64); !ok {
			panic("message_bytes assert float64 failed")
		}
	}

	return &Event{
		Type:         AddNodeEvent,
		IP:           ed.IP,
		Port:         ed.Port,
		ConnectNum:   connNum,
		MessageBytes: msgBytes,
	}
}

func (e *Event) Key() string {
	return fmt.Sprintf("%s:%s", e.IP, e.Port)
}

package source

import (
	"context"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/coderc/im/common/config"
	"github.com/coderc/im/common/discovery"
)

func Init() {
	eventChan = make(chan *Event)
	ctx := context.Background()
	go DataHandler(&ctx)

	if config.IsDebug() {
		ctx := context.Background()
		testServiceRegister(&ctx, "0000", "test_node0")
		testServiceRegister(&ctx, "1111", "test_node1")
		testServiceRegister(&ctx, "2222", "test_node2")
	}
}

// DataHandler 服务发现事件
func DataHandler(ctx *context.Context) {
	dis := discovery.NewServiceDiscovery(ctx)
	defer dis.Close()
	setFunc := func(key, val string) {
		if ed, err := discovery.UnMarshal([]byte(val)); err == nil {
			if event := NewEvent(ed); ed != nil {
				event.Type = AddNodeEvent
				eventChan <- event
			}
		} else {
			logger.CtxErrorf(*ctx, "DataHandler.setFunc.err: %s", err.Error())
		}
	}

	delFunc := func(key, val string) {
		if ed, err := discovery.UnMarshal([]byte(val)); err == nil {
			if event := NewEvent(ed); ed != nil {
				event.Type = DelNodeEvent
				eventChan <- event
			}
		} else {
			logger.CtxErrorf(*ctx, "DataHandler.delFunc.err: %s", err.Error())
		}
	}

	err := dis.WatchService(config.GetServicePathForIPconfig(), setFunc, delFunc)
	if err != nil {
		panic(err)
	}
}

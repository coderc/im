package source

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/coderc/im/common/config"
	"github.com/coderc/im/common/discovery"
)

func testServiceRegister(ctx *context.Context, port, node string) {
	// 模拟服务发现事件 (新节点加入)
	go func() {
		ed := discovery.EndPointInfo{
			IP:   "1.1.1.1",
			Port: port,
			MetaData: map[string]interface{}{
				"connect_num":   float64(rand.Int63n(1 << 60)),
				"message_bytes": float64(rand.Int63n(1 << 50)),
			},
		}

		sr, err := discovery.NewServiceRegister(ctx, fmt.Sprintf("%s/%s", config.GetServicePathForIPconfig(), node), &ed, time.Now().Unix())
		if err != nil {
			panic(err)
		}

		go sr.ListenLeaseRespChan()
		for {
			ed = discovery.EndPointInfo{
				IP:   "127.0.0.1",
				Port: port,
				MetaData: map[string]interface{}{
					"connect_num":   float64(rand.Int63n(1 << 60)),
					"message_bytes": float64(rand.Int63n(1 << 50)),
				},
			}
			sr.UpdateValue(&ed)
			time.Sleep(1 * time.Second)
		}
	}()
}

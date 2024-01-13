package discovery

import (
	"context"
	"sync"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/coderc/im/common/config"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
)

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	cli  *clientv3.Client // etcd client
	lock sync.Mutex
	ctx  *context.Context
}

// NewServiceDiscovery 新建服务发现
func NewServiceDiscovery(ctx *context.Context) *ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEndPointsForDiscovery(), // etcd connect endpoints
		DialTimeout: config.GetTimeoutForDiscovery(),   // etcd connect dial timeout
	})
	if err != nil {
		logger.Fatal(err)
	}
	return &ServiceDiscovery{
		cli: cli,
		ctx: ctx,
	}
}

// WatchService 初始化服务列表和监听
func (s *ServiceDiscovery) WatchService(prefix string, set, del func(key, val string)) error {
	resp, err := s.cli.Get(*s.ctx, prefix, clientv3.WithPrefix()) // 初始化列表
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		set(string(ev.Key), string(ev.Value))
	}

	s.watcher(prefix, resp.Header.Revision+1, set, del) // 启动监听
	return nil
}

// watcher 监听前缀
func (s *ServiceDiscovery) watcher(prefix string, rev int64, set, del func(key, val string)) {
	rch := s.cli.Watch(*s.ctx, prefix, clientv3.WithPrefix(), clientv3.WithRev(rev))
	logger.CtxInfof(*s.ctx, "watching prefix:%s now...", prefix)
	for wResp := range rch {
		for _, ev := range wResp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				set(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				del(string(ev.Kv.Key), string(ev.Kv.Value))
			}
		}
	}
}

func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}

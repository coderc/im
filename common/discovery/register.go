package discovery

import (
	"context"
	"log"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/coderc/im/common/config"
	"go.etcd.io/etcd/clientv3"
)

// ServiceRegister 服务发现
type ServiceRegister struct {
	cli           *clientv3.Client                        // etcd client
	leaseId       clientv3.LeaseID                        // 租约Id
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse // 租约keepalive channel
	key           string
	val           string
	ctx           *context.Context
}

// NewServiceRegister 新增注册服务
func NewServiceRegister(ctx *context.Context, key string, endPointInfo *EndPointInfo, lease int64) (*ServiceRegister, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEndPointsForDiscovery(),
		DialTimeout: config.GetTimeoutForDiscovery(),
	})

	if err != nil {
		log.Fatal(err)
	}

	ser := &ServiceRegister{
		cli: cli,
		key: key,
		val: endPointInfo.Marshal(),
		ctx: ctx,
	}

	// 申请租约&设置keepalive
	if err := ser.putKeyWithLease(lease); err != nil {
		return nil, err
	}
	return ser, nil
}

func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	resp, err := s.cli.Grant(*s.ctx, lease)
	if err != nil {
		return err
	}

	// 注册服务并绑定租约
	_, err = s.cli.Put(*s.ctx, s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	// 设置续租 定期发送需求请求
	leaseRespChan, err := s.cli.KeepAlive(*s.ctx, resp.ID)
	if err != nil {
		return err
	}
	s.leaseId = resp.ID
	s.keepAliveChan = leaseRespChan
	return nil
}

// UpdateValue 更新服务的信息
func (s *ServiceRegister) UpdateValue(val *EndPointInfo) error {
	value := val.Marshal()
	_, err := s.cli.Put(*s.ctx, s.key, value, clientv3.WithLease(s.leaseId))
	if err != nil {
		return err
	}
	s.val = value
	logger.CtxInfof(*s.ctx, "ServiceRegister.updateValue leaseId:%d Put key:%s,val:%s, success!", s.leaseId, s.key, s.val)
	return nil
}

// ListenLeaseRespChan 监听续租情况
func (s *ServiceRegister) ListenLeaseRespChan() {
	for leaseKeepResp := range s.keepAliveChan {
		logger.CtxInfof(*s.ctx, "lease success leaseID:%d, Put key:%s,val:%s reps:+%v",
			s.leaseId, s.key, s.val, leaseKeepResp)
	}
	logger.CtxInfof(*s.ctx, "lease failed !!!  leaseID:%d, Put key:%s,val:%s", s.leaseId, s.key, s.val)
}

func (s *ServiceRegister) Close() error {
	// 撤销租约
	if _, err := s.cli.Revoke(context.Background(), s.leaseId); err != nil {
		return err
	}
	logger.CtxInfof(*s.ctx, "lease close !!! leaseId:%d, Put key:%s,val:%s success", s.leaseId, s.key, s.val)

	return s.cli.Close()
}

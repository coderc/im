package domain

import (
	"sort"
	"sync"

	"github.com/coderc/im/ipconfig/source"
)

type Dispatcher struct {
	candidateTable map[string]*EndPoint
	sync.RWMutex
}

var (
	dp *Dispatcher // 全局服务调度器，用于定制和处理流量分配的逻辑
)

func Init() {
	dp = &Dispatcher{}
	dp.candidateTable = make(map[string]*EndPoint)
	go func() {
		for event := range source.EventChan() {
			switch event.Type {
			case source.AddNodeEvent:
				dp.addNode(event)
			case source.DelNodeEvent:
				dp.delNode(event)
			}
		}
	}()
}

func Dispatch(ctx *IpConfigContext) []*EndPoint {
	// 获取候选endpoint
	eds := dp.getCandidateEndpoint(ctx)
	// 计算得分
	for _, ed := range eds {
		ed.CalculateScore(ctx)
	}
	// 排序
	sort.Slice(eds, func(i, j int) bool {
		// 优先基于active score进行排序，当active score相等时基于static score排序
		return eds[i].ActiveScore > eds[j].ActiveScore || (eds[i].ActiveScore == eds[j].ActiveScore && eds[i].StaticScore > eds[j].StaticScore)
	})
	return eds
}

// getCandidateEndpoint 获取调度器当前持有的下游服务的endpoint list
func (dp *Dispatcher) getCandidateEndpoint(ctx *IpConfigContext) []*EndPoint {
	dp.RLock()
	defer dp.RUnlock()
	candidateList := make([]*EndPoint, 0, len(dp.candidateTable))
	for _, ed := range dp.candidateTable {
		candidateList = append(candidateList, ed)
	}
	return candidateList
}

// delNode 下游服务下线
func (dp *Dispatcher) delNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	delete(dp.candidateTable, event.Key())
}

// addNode 下游服务上线/更新
func (dp *Dispatcher) addNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	var (
		ed *EndPoint
		ok bool
	)
	if ed, ok = dp.candidateTable[event.Key()]; !ok {
		ed = NewEndPoint(event.IP, event.Port)
		dp.candidateTable[event.Key()] = ed
	}
	ed.UpdateStat(&Stat{
		ConnectNum:   event.ConnectNum,
		MessageBytes: event.MessageBytes,
	})
}

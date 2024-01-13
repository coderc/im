package domain

const (
	windowSize = 5
)

// stateWindow 记录和处理某一个服务的多次状态更新，用于计算服务的active score 和 static score
type stateWindow struct {
	statQueue []*Stat
	statChan  chan *Stat
	sumStat   *Stat
	idx       int64
}

func newStateWindow() *stateWindow {
	return &stateWindow{
		statQueue: make([]*Stat, windowSize),
		statChan:  make(chan *Stat),
		sumStat:   &Stat{},
	}
}

func (s *stateWindow) getStat() *Stat {
	res := s.sumStat.Clone()
	res.Avg(min(windowSize, float64(s.idx)))
	return res
}

func (s *stateWindow) appendStat(st *Stat) {
	// 减去即将被删除的state
	s.sumStat.Sub(s.statQueue[s.idx%windowSize])
	// 更新
	s.statQueue[s.idx%windowSize] = st
	// 计算新的sum state
	s.sumStat.Add(st)
	s.idx++
}

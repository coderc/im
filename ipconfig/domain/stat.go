package domain

import "math"

type Stat struct {
	ConnectNum   float64
	MessageBytes float64
}

func (s *Stat) CalculateActiveScore() float64 {
	return getGB(s.MessageBytes)
}

func (s *Stat) CalculateStaticScore() float64 {
	return s.ConnectNum
}

func (s *Stat) Avg(num float64) {
	if num == 0 {
		return
	}
	s.ConnectNum /= num
	s.MessageBytes /= num
}

func (s *Stat) Clone() *Stat {
	newStat := &Stat{
		ConnectNum:   s.ConnectNum,
		MessageBytes: s.MessageBytes,
	}
	return newStat
}

func (s *Stat) Add(st *Stat) {
	if s == nil {
		return
	}
	s.ConnectNum += st.ConnectNum
	s.MessageBytes += st.MessageBytes
}

func (s *Stat) Sub(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum -= st.ConnectNum
	s.MessageBytes -= st.MessageBytes
}

func getGB(m float64) float64 {
	return decimal(m / (1 << 30))
}

func decimal(val float64) float64 {
	return math.Trunc(val*1e2+0.5) * 1e-2
}

func min[T int | int64 | int32 | float32 | float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

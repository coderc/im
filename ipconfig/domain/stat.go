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
	if s == nil {
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

func min(a, b, c float64) float64 {
	m := func(i, j float64) float64 {
		if i < j {
			return i
		}
		return j
	}
	return m(a, m(b, c))
}

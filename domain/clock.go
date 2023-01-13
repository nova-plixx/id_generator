package domain

import "time"

type Clock interface {
	EpochMilli() int64
}

type SystemClock struct{}

func (s *SystemClock) EpochMilli() int64 {
	return time.Now().UnixMilli()
}

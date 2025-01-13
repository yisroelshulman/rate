package main

import (
	"sync"
	"time"
)

type UnbufferedLimiter struct {
	mu         *sync.Mutex
	index      int
	interval   time.Duration
	timeStamps []time.Time
}

func NewUnbufferedLimiter(limit int, interval time.Duration) *UnbufferedLimiter {
	return &UnbufferedLimiter{
		mu:         &sync.Mutex{},
		index:      0,
		interval:   interval,
		timeStamps: make([]time.Time, limit),
	}
}

func (l *UnbufferedLimiter) Wait(timeout *time.Duration) error {
	for {
		remaining, err := l.TryWait()
		if err == nil {
			return nil
		}
		time.Sleep(remaining)
	}
}

func (l *UnbufferedLimiter) TryWait() (time.Duration, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	remaining := l.interval - time.Since(l.timeStamps[l.index])
	if remaining < 0 {
		l.timeStamps[l.index] = time.Now()
		l.index = incrementIndex(l.index, len(l.timeStamps))
		return 0, nil
	}
	return remaining, &LimiterOverLimit{message: "limit reached"}
}

package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Limiter struct {
	mu    *sync.Mutex
	Count int
	Limit int
	Rate  time.Duration
}

func NewLimiter(limit int, rate time.Duration) *Limiter {
	l := &Limiter{
		mu:    &sync.Mutex{},
		Count: 0,
		Limit: limit,
		Rate:  rate,
	}
	go l.resetRate(rate)
	return l
}

func (l *Limiter) resetRate(rate time.Duration) {
	ticker := time.NewTicker(rate)
	for range ticker.C {
		l.reset()
	}
}

func (l *Limiter) reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Printf("reset with val: %d\n", l.Count)
	l.Count = 0
}

func (l *Limiter) Wait(timeout *time.Duration) error {
	if timeout != nil {
		return l.waitTime(*timeout)
	}
	for {
		l.mu.Lock()
		if l.Count < l.Limit {
			l.Count++
			l.mu.Unlock()
			return nil
		}
		l.mu.Unlock()
	}
}

func (l *Limiter) waitTime(timeout time.Duration) error {
	to := false
	go timeOut(&to, timeout)
	for to {
		l.mu.Lock()
		if l.Count < l.Limit {
			l.Count++
			l.mu.Unlock()
			return nil
		}
		l.mu.Unlock()
	}
	return errors.New("limiter timed out")
}

func timeOut(timedOut *bool, to time.Duration) {
	ticker := time.NewTicker(to)
	for range ticker.C {
		*timedOut = false
	}
}

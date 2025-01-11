package main

import (
	"fmt"
	"sync"
	"time"
)

type LimiterFull struct {
	message string
}

func (l *LimiterFull) Error() string {
	return l.message
}

type LimiterTimedOut struct {
	message string
}

func (l *LimiterTimedOut) Error() string {
	return l.message
}

type Limiter struct {
	mu     *sync.Mutex
	Count  int
	Limit  int
	Rate   time.Duration
	Buffer *buffer
}

func NewLimiter(limit, capacity int, rate time.Duration) *Limiter {
	l := &Limiter{
		mu:     &sync.Mutex{},
		Count:  0,
		Limit:  limit,
		Rate:   rate,
		Buffer: newBuffer(capacity),
	}
	l.startLimiter(rate)
	return l
}

func (l *Limiter) startLimiter(rate time.Duration) {
	go l.resetLoop(rate)
	go l.processLoop()
}

func (l *Limiter) resetLoop(rate time.Duration) {
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

func (l *Limiter) processLoop() {
	for {
		l.mu.Lock()
		if l.Count < l.Limit {
			if ok := l.Buffer.remove(); ok {
				l.Count++
			}

		}
		l.mu.Unlock()
	}
}

func (l *Limiter) Wait(timeout *time.Duration) error {
	if timeout != nil {
		return l.waitTime(*timeout)
	}
	waiting := true
	if ok := l.Buffer.add(&waiting); !ok {
		return &LimiterFull{message: "buffer full"}
	}
	for {
		if !waiting {
			return nil
		}
	}
}

func (l *Limiter) waitTime(timeout time.Duration) error {
	run := true
	waiting := true
	if ok := l.Buffer.add(&waiting); !ok {
		return &LimiterFull{message: "buffer full"}
	}
	go timeOut(&run, timeout)
	for run {
		if !waiting {
			return nil
		}
	}
	return &LimiterTimedOut{message: "timed out"}
}

func timeOut(timedOut *bool, to time.Duration) {
	ticker := time.NewTicker(to)
	for range ticker.C {
		*timedOut = false
	}
}

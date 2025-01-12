package main

import (
	"fmt"
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
	limit      int
	rate       time.Duration
	timeStamps []time.Time
	buffer     *buffer
}

func NewLimiter(limit, capacity int, rate time.Duration) *Limiter {
	l := &Limiter{
		limit:      limit,
		rate:       rate,
		timeStamps: make([]time.Time, limit),
		buffer:     newBuffer(capacity),
	}
	go l.processLoop()
	return l
}

func (l *Limiter) processLoop() {
	index := 0
	for {
		if time.Since(l.timeStamps[index]) > l.rate {
			if ok := l.buffer.remove(); !ok {
				continue
			}
			fmt.Println(time.Now())
			l.timeStamps[index] = time.Now()
			index = (index + 1) % l.limit
			continue
		}
		time.Sleep(l.rate - time.Since(l.timeStamps[index]))
	}
}

func (l *Limiter) Wait(timeout *time.Duration) error {
	if timeout != nil {
		return l.waitTime(*timeout)
	}
	waiting := true
	if ok, _ := l.buffer.add(&waiting); !ok {
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
	if ok, _ := l.buffer.add(&waiting); !ok {
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

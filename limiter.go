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

type waitStat struct {
	waiting  bool
	timedOut bool
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
	wait := waitStat{
		waiting: true,
	}
	if ok := l.buffer.add(&wait); !ok {
		return &LimiterFull{message: "buffer full"}
	}
	for {
		if !wait.waiting {
			return nil
		}
	}
}

func (l *Limiter) waitTime(timeout time.Duration) error {
	run := true
	wait := waitStat{
		waiting: true,
	}
	if ok := l.buffer.add(&wait); !ok {
		return &LimiterFull{message: "buffer full"}
	}
	go timeOut(&run, timeout)
	for run {
		if !wait.waiting {
			return nil
		}
	}
	wait.timedOut = true
	l.buffer.timedOutSignal()
	return &LimiterTimedOut{message: "timed out"}
}

func timeOut(timedOut *bool, to time.Duration) {
	ticker := time.NewTicker(to)
	for range ticker.C {
		*timedOut = false
	}
}

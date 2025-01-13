package main

import (
	"fmt"
	"time"
)

type waitSignals struct {
	waiting  bool
	timedOut bool
}

type BufferedLimiter struct {
	limit      int
	rate       time.Duration
	timeStamps []time.Time
	buffer     *buffer
}

func NewBufferedLimiter(limit, capacity int, rate time.Duration) *BufferedLimiter {
	if limit == 0 {
		limit = 1
	}
	if capacity == 0 {
		capacity = 1
	}
	l := &BufferedLimiter{
		limit:      limit,
		rate:       rate,
		timeStamps: make([]time.Time, limit),
		buffer:     newBuffer(capacity),
	}
	go l.processLoop()
	return l
}

func (l *BufferedLimiter) processLoop() {
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

func (l *BufferedLimiter) Wait(timeout *time.Duration) error {
	if timeout != nil {
		return l.waitWithTimeOut(*timeout)
	}
	wait := waitSignals{
		waiting: true,
	}
	if ok := l.buffer.add(&wait); !ok {
		return &LimiterBufferFull{message: "permission denied: buffer full"}
	}
	for {
		if !wait.waiting {
			return nil
		}
	}
}

func (l *BufferedLimiter) waitWithTimeOut(timeout time.Duration) error {
	start := time.Now()
	wait := waitSignals{
		waiting:  true,
		timedOut: false,
	}
	if ok := l.buffer.add(&wait); !ok {
		return &LimiterBufferFull{message: "permission denied: buffer full"}
	}
	for time.Since(start) < timeout {
		if !wait.waiting {
			return nil
		}
	}
	wait.timedOut = true
	l.buffer.timedOutSignal()
	return &LimiterWaitTimedOut{message: "permission denied: timed out"}
}

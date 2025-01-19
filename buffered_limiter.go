package main

import (
	"time"
)

// A BufferedLimiter is a limiter with an internal buffer that limits permissions at a given rate
// and time interval
type BufferedLimiter struct {
	rate       int
	interval   time.Duration
	timeStamps []time.Time
	buffer     *buffer
}

// NewBufferedLimiter returns a new BufferedLimiter given a capacity, rate, and interval.
//
// The BufferedLimiter will limit permissions at the rate per time interval sending additional
// requests to the buffer for later processing. The purpose of the buffer is to prevent starvation
// of requests.
//
// If the received rate <= 0, the default rate of 1 will be used, likewise if a capacity <= to 0 the
// capacity for the buffer will be set to 1. If the interval received is <= 0 the time interval will
// default to 1 millisecond. This is to prevent the BufferedLimiter from erroring during use and to
// ensure the caller gets a working limiter without checking for errors.
func NewBufferedLimiter(rate, capacity int, interval time.Duration) *BufferedLimiter {
	if rate <= 0 {
		rate = 1
	}
	if capacity <= 0 {
		capacity = 1
	}
	if interval <= 0 {
		interval = time.Millisecond
	}
	l := &BufferedLimiter{
		rate:       rate,
		interval:   interval,
		timeStamps: make([]time.Time, rate),
		buffer:     newBuffer(capacity),
	}
	go l.permissionApprovalLoop()
	return l
}

// handles the logic to process permission approvals for requests from the buffer.
// This runs on its own goroutine as it loops indefinitely
func (l *BufferedLimiter) permissionApprovalLoop() {
	index := 0
	for {
		if time.Since(l.timeStamps[index]) > l.interval {
			if ok := l.buffer.remove(); !ok { // buffer empty
				continue
			}
			l.timeStamps[index] = time.Now()
			index = incrementIndex(index, l.rate)
			continue
		}
		time.Sleep(l.interval - time.Since(l.timeStamps[index]))
	}
}

// Wait returns when the limiter grants permission or times out.
//
// The Wait receiver for the BufferedLimiter blocks until the permission is granted or the request
// times out. When permission is granted, a nil value is returned. Otherwise, an error will be
// returned depending on what caused the limiter to deny permission, (LimiterBufferFull,
// LimiterWaitTimedOut).
//
// Calls to Wait are thread safe in the sense that the limiter won't break and will enforce the rate
// per interval regardless of how many threads share the limiter. The BufferedLimiter ensures that
// approvals are granted in the order the limiter received the requests for permission, (with timed
// out requests being ignored), preventing request starvation. This does not guarantee that the
// execution of whatever was limited will be in order only that the permissions are granted in order.
// Of course, there are no guarantees when a LimiterBufferFull error is returned.
func (l *BufferedLimiter) Wait(timeout *time.Duration) error {
	if timeout != nil {
		return l.waitWithTimeOut(*timeout)
	}
	access := permissionStatus{
		granted: false,
	}
	if ok := l.buffer.add(&access); !ok {
		return &LimiterBufferFull{message: "permission denied: buffer full"}
	}
	for {
		if access.granted {
			return nil
		}
	}
}

// handles the timeout logic for the BufferedLimiter sending a signal to the buffer when the request
// timed out so the permissionApprovalLoop knows to ignore it
func (l *BufferedLimiter) waitWithTimeOut(timeout time.Duration) error {
	start := time.Now()
	access := permissionStatus{
		granted:  false,
		timedOut: false,
	}
	if ok := l.buffer.add(&access); !ok {
		return &LimiterBufferFull{message: "permission denied: buffer full"}
	}
	for time.Since(start) < timeout {
		if access.granted {
			return nil
		}
	}
	access.timedOut = true
	l.buffer.timedOutSignal()
	return &LimiterWaitTimedOut{message: "permission denied: timed out"}
}

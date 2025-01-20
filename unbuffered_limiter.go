package rate

import (
	"sync"
	"time"
)

// An UnbufferdLimiter is a limiter without an internal buffer that limits permission at a given
// rate and time interval
type UnbufferedLimiter struct {
	mu         *sync.Mutex
	index      int
	interval   time.Duration
	timeStamps []time.Time
}

// NewUnbufferedLimiter returns a new UnbufferedLimiter given a rate and a time interval.
//
// The Unbufferedlimiter will limit permissions at the provided rate for the given time interval.
//
// If the rate received <= 0 the rate will default to 1 and if the received interval <= 0 it will be
// set to 1 millisecond. This is because the Unbufferedlimiter must have a non-zero rate and
// interval for simplicity and ease of use to prevent the UnbufferedLimiter from erroring during use
// and not have the NewUnbufferedLimiter function return an error.
func NewUnbufferedLimiter(rate int, interval time.Duration) *UnbufferedLimiter {
	if rate <= 0 {
		rate = 1
	}
	if interval <= 0 {
		interval = time.Millisecond
	}
	return &UnbufferedLimiter{
		mu:         &sync.Mutex{},
		index:      0,
		interval:   interval,
		timeStamps: make([]time.Time, rate),
	}
}

// Wait returns when the limiter grants permission or times out.
//
// The Wait receiver on the UnbufferedLimiter blocks until permission is granted or the request
// times out. When permission is granted a nil value is returned and when the request times out a
// LimiterWaitTimedOut error is returned.
//
// Additionally it is thread safe in the sense that the limiter will still enforce the rate
// regardless of how many threads share the same Unbufferedlimiter. However, there is no guarantee
// as to which order the Unbufferedlimiter will grant permission or that permission will ever be
// granted if there are a large number of requests (request starvation).
func (l *UnbufferedLimiter) Wait(timeout *time.Duration) error {
	if timeout != nil {
		return l.waitWithTimeout(*timeout)
	}
	for {
		remaining, err := l.TryWait()
		if err == nil {
			return nil
		}
		time.Sleep(remaining)
	}
}

// handles the timeout logic for the wait with timeout
func (l *UnbufferedLimiter) waitWithTimeout(timeout time.Duration) error {
	start := time.Now()
	for time.Since(start) < timeout {
		remaining, err := l.TryWait()
		if err == nil {
			return nil
		}
		if (time.Since(start) + remaining) > timeout {
			time.Sleep(timeout - time.Since(start))
			break
		}
		time.Sleep(remaining)
	}
	return &LimiterWaitTimedOutError{message: "permission denied: timed out"}
}

// TryWait returns whether or not the Unbufferedlimiter granted permission.
//
// The TryWait receiver is non-blocking and returns immediately. If permission is granted the error
// is nil. If permission is not granted a LimiterOverLimit error is returned and the time duration
// until the next potential approval can occur.
//
// TryWait is a threadsafe function allowing multiple threads to share the same Unbufferedlimiter.
func (l *UnbufferedLimiter) TryWait() (time.Duration, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	remaining := l.interval - time.Since(l.timeStamps[l.index])
	if remaining < 0 {
		l.timeStamps[l.index] = time.Now()
		l.index = incrementIndex(l.index, len(l.timeStamps))
		return 0, nil
	}
	return remaining, &LimiterOverLimitError{message: "permission denied: limit reached"}
}

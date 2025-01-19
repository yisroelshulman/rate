package rate

// LimiterWaitTimedOut is the error returned when limiter.Wait times out
type LimiterWaitTimedOut struct {
	message string
}

func (l *LimiterWaitTimedOut) Error() string {
	return l.message
}

// LimiterOverLimit is the error returned when the unbufferedlimiter.TryWait fails
type LimiterOverLimit struct {
	message string
}

func (l *LimiterOverLimit) Error() string {
	return l.message
}

// LimiterBufferFull is the error returned when the buffer of the bufferedlimiter is full
type LimiterBufferFull struct {
	message string
}

func (l *LimiterBufferFull) Error() string {
	return l.message
}

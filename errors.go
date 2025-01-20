package rate

// LimiterWaitTimedOutError is the error returned when limiter.Wait times out
type LimiterWaitTimedOutError struct {
	message string
}

func (l *LimiterWaitTimedOutError) Error() string {
	return l.message
}

// LimiterOverLimitError is the error returned when the unbufferedlimiter.TryWait fails
type LimiterOverLimitError struct {
	message string
}

func (l *LimiterOverLimitError) Error() string {
	return l.message
}

// LimiterBufferFullError is the error returned when the buffer of the bufferedlimiter is full
type LimiterBufferFullError struct {
	message string
}

func (l *LimiterBufferFullError) Error() string {
	return l.message
}

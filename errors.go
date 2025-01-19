package rate

type LimiterWaitTimedOut struct {
	message string
}

func (l *LimiterWaitTimedOut) Error() string {
	return l.message
}

type LimiterOverLimit struct {
	message string
}

func (l *LimiterOverLimit) Error() string {
	return l.message
}

type LimiterBufferFull struct {
	message string
}

func (l *LimiterBufferFull) Error() string {
	return l.message
}

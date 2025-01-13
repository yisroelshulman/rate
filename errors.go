package main

type LimiterOverLimit struct {
	message string
}

func (l *LimiterOverLimit) Error() string {
	return l.message
}

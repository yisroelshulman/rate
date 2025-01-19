package main

import (
	"sync"
	"testing"
	"time"
)

func TestBufferedWait(t *testing.T) {
	tests := []struct {
		name        string
		limiter     *BufferedLimiter
		sleep       time.Duration
		count       int
		minDuration time.Duration
		maxDuration time.Duration
		wantErr     []bool
	}{
		{
			name:        "Limiter within limit",
			limiter:     NewBufferedLimiter(5, 5, time.Second),
			count:       5,
			sleep:       time.Second / 10,
			minDuration: 0,
			maxDuration: time.Second,
			wantErr:     []bool{false, false, false, false, false},
		},
		{
			name:        "Limiter buffer full",
			limiter:     NewBufferedLimiter(2, 2, time.Second),
			count:       5,
			sleep:       time.Second / 10,
			minDuration: time.Second,
			maxDuration: time.Second * 2,
			wantErr:     []bool{false, false, false, false, true},
		},
		{
			name:        "Rate of 0",
			limiter:     NewBufferedLimiter(0, 1, time.Second),
			count:       3,
			sleep:       time.Second / 10,
			minDuration: time.Second,
			maxDuration: time.Second * 2,
			wantErr:     []bool{false, false, true},
		},
		{
			name:        "Capacity of 0",
			limiter:     NewBufferedLimiter(1, 0, time.Second),
			count:       3,
			sleep:       time.Second / 10,
			minDuration: time.Second,
			maxDuration: time.Second * 2,
			wantErr:     []bool{false, false, true},
		},
		{
			name:        "Interval of 0",
			limiter:     NewBufferedLimiter(1, 1, 0),
			count:       3,
			sleep:       time.Millisecond / 80,
			minDuration: time.Millisecond,
			maxDuration: time.Millisecond * 2,
			wantErr:     []bool{false, false, true},
		},
		{
			name:        "Buffer full remove some and add more",
			limiter:     NewBufferedLimiter(1, 2, time.Second),
			count:       6,
			sleep:       time.Second / 5,
			minDuration: time.Second * 3,
			maxDuration: time.Second * 4,
			wantErr:     []bool{false, false, false, true, true, false},
		},
		{
			name:        "High rate: 1000 Request at rate of 500 per second",
			limiter:     NewBufferedLimiter(500, 1000, time.Second),
			count:       1000,
			sleep:       0,
			minDuration: time.Second,
			maxDuration: time.Second * 2,
			wantErr:     []bool{false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := make([]error, tt.count)
			var wg sync.WaitGroup
			start := time.Now()
			for i := 0; i < tt.count; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					gotErr[i] = tt.limiter.Wait(nil)
				}()
				time.Sleep(tt.sleep)
			}
			wg.Wait()
			gotDuration := time.Since(start)

			for i := 0; i < tt.count; i++ {
				remain := i % len(tt.wantErr)
				if (gotErr[i] != nil) != tt.wantErr[remain] {
					t.Errorf("BufferedLimiter.Wait() test %v/%v, wantErr %v, got %v", i+1, tt.count, tt.wantErr[i], gotErr[i])
				}
			}

			if gotDuration < tt.minDuration || gotDuration > tt.maxDuration {
				t.Errorf("BufferedLimiter.Wait(), wantDuration in range (%v, %v), got %v", tt.minDuration, tt.maxDuration, gotDuration)
			}
		})
	}
}

func TestBufferedWaitWithTimeout(t *testing.T) {
	tests := []struct {
		name     string
		limiter  *BufferedLimiter
		count    int
		timeouts []time.Duration
		wantErr  []bool
	}{
		{
			name:     "Single timeout",
			limiter:  NewBufferedLimiter(1, 100, time.Second),
			count:    2,
			timeouts: []time.Duration{0, time.Second / 2},
			wantErr:  []bool{false, true},
		},
		{
			name:     "Multiple timeouts with not timing out at the end",
			limiter:  NewBufferedLimiter(1, 100, time.Second),
			count:    6,
			timeouts: []time.Duration{0, 0, time.Second / 2, time.Second / 2, time.Second * 3, 0},
			wantErr:  []bool{false, false, true, true, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			gotErr := make([]error, tt.count)

			for i := 0; i < tt.count; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					to := &tt.timeouts[i]
					if *to == 0 {
						to = nil
					}
					gotErr[i] = tt.limiter.Wait(to)
				}()
				time.Sleep(time.Millisecond)
			}
			wg.Wait()

			for i := 0; i < tt.count; i++ {
				if (gotErr[i] != nil) != tt.wantErr[i] {
					t.Errorf("UnbufferedLimiter.Wait(with timeout), wantErr %v, got %v", tt.wantErr[i], gotErr[i])
				}
			}
		})
	}
}

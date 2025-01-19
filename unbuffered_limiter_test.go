package rate

import (
	"testing"
	"time"
)

func TestUnbufferedWait(t *testing.T) {
	tests := []struct {
		name        string
		limiter     *UnbufferedLimiter
		count       int
		minDuration time.Duration
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name:        "Limit 5 requests at 3 per second",
			limiter:     NewUnbufferedLimiter(3, time.Second),
			count:       5,
			minDuration: time.Second,
			maxDuration: time.Second * 2,
			wantErr:     false,
		},
		{
			name:        "Limit 31 requests at 10 per second",
			limiter:     NewUnbufferedLimiter(10, time.Second),
			count:       31,
			minDuration: time.Second * 3,
			maxDuration: time.Second * 4,
			wantErr:     false,
		},
		{
			name:        "Rate of 0",
			limiter:     NewUnbufferedLimiter(0, time.Second),
			count:       3,
			minDuration: time.Second * 2,
			maxDuration: time.Second * 3,
			wantErr:     false,
		},
		{
			name:        "Interval of 0",
			limiter:     NewUnbufferedLimiter(1, 0),
			count:       5,
			minDuration: time.Millisecond * 4,
			maxDuration: time.Millisecond * 5,
			wantErr:     false,
		},
		{
			name:        "500 requests in 1 second",
			limiter:     NewUnbufferedLimiter(500, time.Second),
			count:       500,
			minDuration: 0,
			maxDuration: time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			for i := 0; i < tt.count; i++ {
				gotErr := tt.limiter.Wait(nil)
				if (gotErr != nil) != tt.wantErr {
					t.Errorf("UnbufferedLimiter.Wait(), wantErr %v, got %v", tt.wantErr, gotErr)
				}
			}
			gotDuration := time.Since(start)
			if gotDuration < tt.minDuration || gotDuration > tt.maxDuration {
				t.Errorf("UnbufferedLimiter.Wait(), wantDuration in range (%v, %v) got %v", tt.minDuration, tt.maxDuration, gotDuration)
			}
		})
	}
}

func TestUnbufferedWaitWithTimeout(t *testing.T) {
	tests := []struct {
		name     string
		limiter  *UnbufferedLimiter
		count    int
		timeouts []time.Duration
		wantErr  []bool
	}{
		{
			name:     "Single time out",
			limiter:  NewUnbufferedLimiter(1, time.Second*2),
			count:    2,
			timeouts: []time.Duration{0, time.Second * 1},
			wantErr:  []bool{false, true},
		},
		{
			name:     "Multiple time outs with no time out at the end",
			limiter:  NewUnbufferedLimiter(3, time.Second*3),
			count:    6,
			timeouts: []time.Duration{0, 0, 0, time.Second, time.Second, time.Second * 2},
			wantErr:  []bool{false, false, false, true, true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.count; i++ {
				to := &tt.timeouts[i]
				if *to == 0 {
					to = nil
				}
				gotErr := tt.limiter.Wait(to)
				if (gotErr != nil) != tt.wantErr[i] {
					t.Errorf("UnbufferedLimiter.Wait(with timeout) %v/%v, wantErr %v, got %v", i+1, tt.count, tt.wantErr[i], gotErr)
				}
			}
		})
	}
}

func TestUnbufferedTryWait(t *testing.T) {
	tests := []struct {
		name         string
		limiter      *UnbufferedLimiter
		count        int
		sleep        time.Duration
		wantDuration []bool
		wantErr      []bool
	}{
		{
			name:         "5 requests limit 3 Paused for correct rate",
			limiter:      NewUnbufferedLimiter(3, time.Second),
			count:        5,
			sleep:        time.Second / 2,
			wantDuration: []bool{false, false, false, false, false},
			wantErr:      []bool{false, false, false, false, false},
		},
		{
			name:         "2 Requsts Limit 1 no pause",
			limiter:      NewUnbufferedLimiter(1, time.Second),
			count:        2,
			sleep:        0,
			wantDuration: []bool{false, true},
			wantErr:      []bool{false, true},
		},
		{
			name:         "15 requests limit 5, Paused so the middle 5 are overlimit",
			limiter:      NewUnbufferedLimiter(5, time.Second),
			count:        15,
			sleep:        time.Second / 10,
			wantDuration: []bool{false, false, false, false, false, true, true, true, true, true, false, false, false, false, false},
			wantErr:      []bool{false, false, false, false, false, true, true, true, true, true, false, false, false, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.count; i++ {
				gotDuration, gotErr := tt.limiter.TryWait()
				time.Sleep(tt.sleep)

				if (gotErr != nil) != tt.wantErr[i] {
					t.Errorf("UnbufferedLimiter.TryWait(), wantErr %v, got %v", tt.wantErr[i], gotErr)
				}
				if (gotDuration > 0) != tt.wantDuration[i] {
					t.Errorf("UnbufferdLimiter.TryWait(), wantDuration %v, got %v", tt.wantDuration[i], gotDuration)
				}

			}
		})
	}
}

package main

import (
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name        string
		buf         *buffer
		insertions  int
		permissions permissionStatus
		wantOk      []bool
	}{
		{
			name:       "Add one",
			buf:        newBuffer(1),
			insertions: 1,
			permissions: permissionStatus{
				granted:  false,
				timedOut: false,
			},
			wantOk: []bool{true},
		},
		{
			name:       "Add multiple",
			buf:        newBuffer(3),
			insertions: 3,
			permissions: permissionStatus{
				granted:  false,
				timedOut: false,
			},
			wantOk: []bool{true, true, true},
		},
		{
			name:       "Add to full buffer",
			buf:        newBuffer(1),
			insertions: 2,
			permissions: permissionStatus{
				granted:  false,
				timedOut: false,
			},
			wantOk: []bool{true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.insertions; i++ {
				gotOk := tt.buf.add(&tt.permissions)
				if gotOk != tt.wantOk[i] {
					t.Errorf("buffer.add(), wantOk %v, got %v", tt.wantOk[i], gotOk)
				}

			}
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name            string
		buf             *buffer
		insertions      int
		permissions     permissionStatus
		timedOutIndices []int
		removals        int
		wantOk          []bool
	}{
		{
			name:       "Remove one",
			buf:        newBuffer(3),
			insertions: 3,
			permissions: permissionStatus{
				granted:  false,
				timedOut: false,
			},
			timedOutIndices: nil,
			removals:        1,
			wantOk:          []bool{true},
		},
		{
			name:       "Remove multiple",
			buf:        newBuffer(3),
			insertions: 3,
			permissions: permissionStatus{
				granted:  false,
				timedOut: false,
			},
			timedOutIndices: nil,
			removals:        3,
			wantOk:          []bool{true, true, true},
		},
		{
			name:       "Remove from empty",
			buf:        newBuffer(3),
			insertions: 3,
			permissions: permissionStatus{
				granted:  false,
				timedOut: false,
			},
			timedOutIndices: nil,
			removals:        4,
			wantOk:          []bool{true, true, true, false},
		},
		{
			name:       "Remove through single timed out",
			buf:        newBuffer(3),
			insertions: 3,
			permissions: permissionStatus{
				granted:  false,
				timedOut: false,
			},
			timedOutIndices: []int{1},
			removals:        2,
			wantOk:          []bool{true, true},
		},
		{
			name:       "Remove through multiple timed out",
			buf:        newBuffer(4),
			insertions: 4,
			permissions: permissionStatus{
				granted:  false,
				timedOut: false,
			},
			timedOutIndices: []int{1, 2},
			removals:        2,
			wantOk:          []bool{true, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.insertions; i++ {
				tt.buf.add(&tt.permissions)
			}

			for _, i := range tt.timedOutIndices {
				tt.buf.buffer[i] = &permissionStatus{
					granted:  false,
					timedOut: true,
				}
				tt.buf.timedOutSignal()
			}

			for i := 0; i < tt.removals; i++ {
				gotOk := tt.buf.remove()
				if gotOk != tt.wantOk[i] {
					t.Errorf("buffer.remove(), wantOk %v, got %v", tt.wantOk[i], gotOk)
				}
			}

			if !tt.permissions.granted {
				t.Errorf("buffer.remove(), expext permissions.granted ture, got %v", tt.permissions.granted)
			}
		})
	}
}

func TestCleanBuffer(t *testing.T) {
	tests := []struct {
		name            string
		buf             *buffer
		insertions      int
		permissions     permissionStatus
		timedOutIndices []int
		fullInsertions  int
		wantOk          []bool
	}{
		{
			name:       "Adding to full buffer",
			buf:        newBuffer(3),
			insertions: 3,
			permissions: permissionStatus{
				granted:  false,
				timedOut: false,
			},
			timedOutIndices: nil,
			fullInsertions:  1,
			wantOk:          []bool{false},
		},
		{
			name:       "Adding to buffer requiring cleaning with one timed out",
			buf:        newBuffer(3),
			insertions: 3,
			permissions: permissionStatus{
				granted:  false,
				timedOut: false,
			},
			timedOutIndices: []int{1},
			fullInsertions:  2,
			wantOk:          []bool{true, false},
		},
		{
			name:       "Adding to buffer requiring cleaning with consecutive timed out",
			buf:        newBuffer(8),
			insertions: 8,
			permissions: permissionStatus{
				granted:  false,
				timedOut: false,
			},
			timedOutIndices: []int{1, 2, 4, 6},
			fullInsertions:  5,
			wantOk:          []bool{true, true, true, true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.insertions; i++ {
				tt.buf.add(&tt.permissions)
			}

			for _, i := range tt.timedOutIndices {
				tt.buf.buffer[i] = &permissionStatus{
					granted:  false,
					timedOut: true,
				}
				tt.buf.timedOutSignal()
			}

			for i := 0; i < tt.fullInsertions; i++ {
				gotOk := tt.buf.add(&tt.permissions)
				if gotOk != tt.wantOk[i] {
					t.Errorf("buffer.cleanBuffer(), wantOk %v, got %v", tt.wantOk[i], gotOk)
				}
			}

		})
	}
}

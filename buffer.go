package main

import (
	"sync"
)

type buffer struct {
	mu       *sync.Mutex
	buffer   []*bool
	capacity int
	size     int
	insertAt int
	removeAt int
}

func newBuffer(capacity int) *buffer {
	b := &buffer{
		mu:       &sync.Mutex{},
		buffer:   make([]*bool, capacity),
		capacity: capacity,
		size:     0,
		insertAt: 0,
		removeAt: 0,
	}
	return b
}

func (b *buffer) add(lock *bool) (bool, int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.size == b.capacity {
		return false, -1
	}
	b.buffer[b.insertAt] = lock
	b.insertAt = (b.insertAt + 1) % b.capacity
	b.size++
	return true, b.insertAt - 1
}

func (b *buffer) remove() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.size == 0 {
		return false
	}
	*b.buffer[b.removeAt] = false
	b.buffer[b.removeAt] = nil
	b.removeAt = (b.removeAt + 1) % b.capacity
	b.size--
	return true
}

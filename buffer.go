package main

import (
	"fmt"
	"sync"
)

type buffer struct {
	mu       *sync.Mutex
	buffer   []*waitStat
	capacity int
	size     int
	insertAt int
	removeAt int
}

func newBuffer(capacity int) *buffer {
	b := &buffer{
		mu:       &sync.Mutex{},
		buffer:   make([]*waitStat, capacity),
		capacity: capacity,
		size:     0,
		insertAt: 0,
		removeAt: 0,
	}
	return b
}

func (b *buffer) add(lock *waitStat) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.size == b.capacity {
		return false
	}
	b.buffer[b.insertAt] = lock
	b.insertAt = incrementIndex(b.insertAt, b.capacity)
	b.size++
	if b.insertAt == b.removeAt && b.size != b.capacity {
		b.cleanBuffer()
	}
	return true
}

func (b *buffer) cleanBuffer() {
	fmt.Println(">>>>> cleaning buffer")
	buf := make([]*waitStat, b.capacity)
	pos := b.removeAt
	insert := 0
	for i := 0; i < b.capacity; i++ {
		if !b.buffer[pos].timedOut {
			buf[insert] = b.buffer[pos]
			insert++
		}
		pos = incrementIndex(pos, b.capacity)
	}
	b.buffer = buf
	b.insertAt = insert
	b.removeAt = 0
}

func (b *buffer) remove() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.size == 0 {
		return false
	}
	for b.buffer[b.removeAt] == nil || b.buffer[b.removeAt].timedOut {
		b.buffer[b.removeAt] = nil
		b.removeAt = incrementIndex(b.removeAt, b.capacity)
	}
	b.buffer[b.removeAt].waiting = false
	b.buffer[b.removeAt] = nil
	b.removeAt = incrementIndex(b.removeAt, b.capacity)
	b.size--
	return true
}

func (b *buffer) timedOutSignal() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.size--
}

func incrementIndex(index, capacity int) int {
	return (index + 1) % capacity
}

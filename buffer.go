package rate

import (
	"sync"
)

// permissionStatus is used to send signals between the limiter and the buffer.
//
// The limiter signals when the request times out and the buffer signals when permission is granted.
type permissionStatus struct {
	granted  bool
	timedOut bool
}

// a buffer for the BufferedLimiter to keep track of the requests waiting for approval.
type buffer struct {
	mu       *sync.Mutex
	capacity int
	size     int
	insertAt int
	removeAt int
	buffer   []*permissionStatus
}

// returns a new buffer
//
// no need to check if capacity is 0 since this is for internal use only.
func newBuffer(capacity int) *buffer {
	return &buffer{
		mu:       &sync.Mutex{},
		capacity: capacity,
		size:     0,
		insertAt: 0,
		removeAt: 0,
		buffer:   make([]*permissionStatus, capacity),
	}
}

// add the request to the buffer
// returns true if the request was added and false if the buffer is full
func (b *buffer) add(access *permissionStatus) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.size == b.capacity {
		return false
	}
	if b.insertAt == b.removeAt && b.size > 0 { // some requests timed out so cleaning is necessary
		b.cleanBuffer()
	}
	b.buffer[b.insertAt] = access
	b.insertAt = incrementIndex(b.insertAt, b.capacity)
	b.size++
	return true
}

// removes the requests that timed out preserving the order of the requests
func (b *buffer) cleanBuffer() {
	buf := make([]*permissionStatus, b.capacity)
	pos := b.removeAt
	insert := 0
	for i := 0; i < b.capacity; i++ {
		if b.buffer[pos] != nil && !b.buffer[pos].timedOut { // check that the request hasn't timed out
			buf[insert] = b.buffer[pos]
			insert++
		}
		pos = incrementIndex(pos, b.capacity)
	}
	b.buffer = buf
	b.insertAt = insert
	b.removeAt = 0
}

// remove signals to the next requester that access was granted if there is one waiting
//
// returns true if a requester is waiting false otherwise
func (b *buffer) remove() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.size == 0 {
		return false
	}
	for b.buffer[b.removeAt].timedOut { // advance to the next requester
		b.buffer[b.removeAt] = nil
		b.removeAt = incrementIndex(b.removeAt, b.capacity)
	}
	b.buffer[b.removeAt].granted = true
	b.buffer[b.removeAt] = nil
	b.removeAt = incrementIndex(b.removeAt, b.capacity)
	b.size--
	return true
}

// the requester signals to the buffer that the request timed out
func (b *buffer) timedOutSignal() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.size--
}

// increments an array index and wraps around to prevent index out of bounds
func incrementIndex(index, capacity int) int {
	return (index + 1) % capacity
}

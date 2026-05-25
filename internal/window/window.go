// Package window provides a sliding-window line buffer that retains the
// last N lines seen, useful for context extraction around matched events.
package window

import "errors"

// Buffer is a fixed-capacity circular buffer of strings.
type Buffer struct {
	buf  []string
	head int
	size int
	cap  int
}

// New creates a Buffer that retains at most n lines.
// n must be >= 1.
func New(n int) (*Buffer, error) {
	if n < 1 {
		return nil, errors.New("window: capacity must be >= 1")
	}
	return &Buffer{
		buf: make([]string, n),
		cap: n,
	}, nil
}

// Push adds a line to the buffer, evicting the oldest entry when full.
func (b *Buffer) Push(line string) {
	b.buf[b.head] = line
	b.head = (b.head + 1) % b.cap
	if b.size < b.cap {
		b.size++
	}
}

// Lines returns the buffered lines in chronological order (oldest first).
func (b *Buffer) Lines() []string {
	if b.size == 0 {
		return nil
	}
	out := make([]string, b.size)
	start := (b.head - b.size + b.cap) % b.cap
	for i := 0; i < b.size; i++ {
		out[i] = b.buf[(start+i)%b.cap]
	}
	return out
}

// Len returns the number of lines currently held.
func (b *Buffer) Len() int { return b.size }

// Reset clears the buffer without reallocating.
func (b *Buffer) Reset() {
	b.head = 0
	b.size = 0
}

// Package linecount provides utilities for counting lines in a file
// or an io.Reader without fully loading the content into memory.
package linecount

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Counter counts newline-terminated lines from a source.
type Counter struct {
	bufSize int
}

// Option is a functional option for Counter.
type Option func(*Counter)

// WithBufSize sets the internal scanner buffer size in bytes.
func WithBufSize(n int) Option {
	return func(c *Counter) {
		c.bufSize = n
	}
}

// New creates a Counter with optional configuration.
func New(opts ...Option) (*Counter, error) {
	c := &Counter{bufSize: 64 * 1024}
	for _, o := range opts {
		o(c)
	}
	if c.bufSize <= 0 {
		return nil, fmt.Errorf("linecount: bufSize must be positive, got %d", c.bufSize)
	}
	return c, nil
}

// CountReader counts lines from r until EOF or an error.
func (c *Counter) CountReader(r io.Reader) (int64, error) {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, c.bufSize)
	scanner.Buffer(buf, c.bufSize)

	var n int64
	for scanner.Scan() {
		n++
	}
	if err := scanner.Err(); err != nil {
		return n, fmt.Errorf("linecount: scan error: %w", err)
	}
	return n, nil
}

// CountFile counts lines in the file at path.
func (c *Counter) CountFile(path string) (int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("linecount: open %q: %w", path, err)
	}
	defer f.Close()
	return c.CountReader(f)
}

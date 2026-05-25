// Package head provides a reader that returns only the first N lines
// of a log stream, analogous to the Unix `head` command.
package head

import (
	"bufio"
	"errors"
	"io"
)

// Reader reads up to N lines from the underlying io.Reader.
type Reader struct {
	scanner *bufio.Scanner
	max     int
	seen    int
}

// New creates a Reader that will return at most n lines.
// n must be greater than zero.
func New(r io.Reader, n int) (*Reader, error) {
	if n <= 0 {
		return nil, errors.New("head: n must be greater than zero")
	}
	if r == nil {
		return nil, errors.New("head: reader must not be nil")
	}
	return &Reader{
		scanner: bufio.NewScanner(r),
		max:     n,
	}, nil
}

// ReadLine returns the next line and true while lines remain within the
// first n lines. When the limit is reached or the underlying reader is
// exhausted it returns ("", false). Errors from the underlying scanner
// are silently treated as EOF.
func (h *Reader) ReadLine() (string, bool) {
	if h.seen >= h.max {
		return "", false
	}
	if !h.scanner.Scan() {
		return "", false
	}
	h.seen++
	return h.scanner.Text(), true
}

// Lines collects all lines up to n into a slice.
func (h *Reader) Lines() []string {
	var out []string
	for {
		line, ok := h.ReadLine()
		if !ok {
			break
		}
		out = append(out, line)
	}
	return out
}

// Seen returns the number of lines read so far.
func (h *Reader) Seen() int {
	return h.seen
}

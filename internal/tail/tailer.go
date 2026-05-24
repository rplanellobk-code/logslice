package tail

import (
	"bufio"
	"errors"
	"io"
)

// Tailer returns the last N lines from a reader efficiently by maintaining
// a circular buffer of lines.
type Tailer struct {
	n int
}

// New creates a Tailer that will return at most n lines from the end of input.
// Returns an error if n is less than 1.
func New(n int) (*Tailer, error) {
	if n < 1 {
		return nil, errors.New("tail: n must be at least 1")
	}
	return &Tailer{n: n}, nil
}

// Read consumes all lines from r and returns the last n lines.
// The returned slice contains at most n elements, in original order.
func (t *Tailer) Read(r io.Reader) ([]string, error) {
	if r == nil {
		return nil, errors.New("tail: reader must not be nil")
	}

	buf := make([]string, t.n)
	pos := 0
	count := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		buf[pos%t.n] = scanner.Text()
		pos++
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if count == 0 {
		return []string{}, nil
	}

	size := count
	if size > t.n {
		size = t.n
	}

	result := make([]string, size)
	start := pos % t.n
	for i := 0; i < size; i++ {
		result[i] = buf[(start+i)%t.n]
	}
	return result, nil
}

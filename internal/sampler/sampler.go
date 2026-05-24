// Package sampler provides log line sampling strategies for reducing
// output volume while preserving representative data across large archives.
package sampler

import (
	"errors"
	"sync/atomic"
)

// Line represents a single parsed log entry passed through the sampler.
type Line struct {
	Raw       string
	Timestamp int64
}

// Sampler decides whether a given log line should be kept.
type Sampler struct {
	n       uint64 // keep every nth line
	counter atomic.Uint64
}

// New creates a Sampler that retains every nth line.
// n must be >= 1; n=1 keeps every line (no-op sampling).
func New(n uint64) (*Sampler, error) {
	if n == 0 {
		return nil, errors.New("sampler: n must be >= 1")
	}
	return &Sampler{n: n}, nil
}

// Keep returns true if the line should be included in the output.
// It is safe to call from multiple goroutines.
func (s *Sampler) Keep(_ Line) bool {
	v := s.counter.Add(1)
	return v%s.n == 1
}

// Reset resets the internal counter, restarting the sampling window.
func (s *Sampler) Reset() {
	s.counter.Store(0)
}

// Rate returns the configured sampling rate (1-in-n).
func (s *Sampler) Rate() uint64 {
	return s.n
}

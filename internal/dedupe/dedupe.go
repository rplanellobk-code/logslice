// Package dedupe provides line-level deduplication for log streams.
// It tracks a sliding window of recently seen lines and filters out
// exact duplicates within that window.
package dedupe

import "errors"

// Filter holds state for deduplication across a stream of log lines.
type Filter struct {
	seen    map[string]struct{}
	window  []string
	capacity int
}

// New creates a Filter with the given window capacity.
// capacity is the maximum number of unique lines remembered at once.
// Returns an error if capacity is less than 1.
func New(capacity int) (*Filter, error) {
	if capacity < 1 {
		return nil, errors.New("dedupe: capacity must be at least 1")
	}
	return &Filter{
		seen:     make(map[string]struct{}, capacity),
		window:   make([]string, 0, capacity),
		capacity: capacity,
	}, nil
}

// IsDuplicate reports whether line has been seen within the current window.
// If it is new, the line is recorded and, if the window is full, the oldest
// entry is evicted.
func (f *Filter) IsDuplicate(line string) bool {
	if _, ok := f.seen[line]; ok {
		return true
	}
	// Evict oldest entry when at capacity.
	if len(f.window) == f.capacity {
		oldest := f.window[0]
		f.window = f.window[1:]
		delete(f.seen, oldest)
	}
	f.seen[line] = struct{}{}
	f.window = append(f.window, line)
	return false
}

// Reset clears all remembered lines.
func (f *Filter) Reset() {
	f.seen = make(map[string]struct{}, f.capacity)
	f.window = f.window[:0]
}

// Len returns the number of lines currently in the window.
func (f *Filter) Len() int {
	return len(f.window)
}

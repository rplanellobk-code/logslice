// Package burst provides detection of log line bursts — sudden spikes
// in log volume within a rolling time window.
package burst

import (
	"errors"
	"sync"
	"time"
)

// Detector tracks log line arrival times and reports whether the current
// rate exceeds a configured threshold within a sliding window.
type Detector struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	times     []time.Time
}

// New returns a Detector that fires when more than threshold lines arrive
// within window. Both values must be positive.
func New(window time.Duration, threshold int) (*Detector, error) {
	if window <= 0 {
		return nil, errors.New("burst: window must be positive")
	}
	if threshold <= 0 {
		return nil, errors.New("burst: threshold must be positive")
	}
	return &Detector{
		window:    window,
		threshold: threshold,
		times:     make([]time.Time, 0, threshold),
	}, nil
}

// Record registers a new line arrival at now and returns true when the
// number of arrivals within the detector's window exceeds the threshold.
func (d *Detector) Record(now time.Time) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	cutoff := now.Add(-d.window)
	d.evict(cutoff)
	d.times = append(d.times, now)
	return len(d.times) > d.threshold
}

// Count returns the number of arrivals currently within the window.
func (d *Detector) Count(now time.Time) int {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.evict(now.Add(-d.window))
	return len(d.times)
}

// Reset clears all recorded arrivals.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.times = d.times[:0]
}

// evict removes entries older than cutoff. Must be called with mu held.
func (d *Detector) evict(cutoff time.Time) {
	i := 0
	for i < len(d.times) && d.times[i].Before(cutoff) {
		i++
	}
	d.times = d.times[i:]
}

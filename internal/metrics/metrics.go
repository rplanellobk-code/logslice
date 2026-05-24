// Package metrics provides lightweight run-time counters for a logslice
// processing pipeline. It tracks lines read, matched, written, and skipped
// so callers can surface a final summary without pulling in a heavy
// observability library.
package metrics

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

// Counters holds all pipeline metrics. Fields are updated atomically so the
// struct is safe to share across goroutines.
type Counters struct {
	LinesRead    atomic.Int64
	LinesMatched atomic.Int64
	LinesWritten atomic.Int64
	LinesSkipped atomic.Int64
	BytesWritten atomic.Int64
	start        time.Time
}

// New returns a Counters instance with the internal clock started.
func New() *Counters {
	return &Counters{start: time.Now()}
}

// Elapsed returns the duration since New was called.
func (c *Counters) Elapsed() time.Duration {
	return time.Since(c.start)
}

// Summary writes a human-readable summary of all counters to w.
func (c *Counters) Summary(w io.Writer) {
	elapsed := c.Elapsed()
	fmt.Fprintf(w, "elapsed:       %s\n", elapsed.Round(time.Millisecond))
	fmt.Fprintf(w, "lines read:    %d\n", c.LinesRead.Load())
	fmt.Fprintf(w, "lines matched: %d\n", c.LinesMatched.Load())
	fmt.Fprintf(w, "lines written: %d\n", c.LinesWritten.Load())
	fmt.Fprintf(w, "lines skipped: %d\n", c.LinesSkipped.Load())
	fmt.Fprintf(w, "bytes written: %d\n", c.BytesWritten.Load())
}

// Rate returns the number of lines read per second, or 0 if elapsed is 0.
func (c *Counters) Rate() float64 {
	secs := c.Elapsed().Seconds()
	if secs == 0 {
		return 0
	}
	return float64(c.LinesRead.Load()) / secs
}

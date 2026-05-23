package progress

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

// Tracker reports processing progress to a writer at regular intervals.
type Tracker struct {
	out      io.Writer
	ticker   *time.Ticker
	done     chan struct{}
	bytesIn  atomic.Int64
	linesIn  atomic.Int64
	linesOut atomic.Int64
}

// NewTracker creates a Tracker that writes progress to out every interval.
// Call Stop when processing is complete.
func NewTracker(out io.Writer, interval time.Duration) *Tracker {
	t := &Tracker{
		out:  out,
		done: make(chan struct{}),
	}
	if interval > 0 && out != nil {
		t.ticker = time.NewTicker(interval)
		go t.loop()
	}
	return t
}

// AddBytes records n bytes read from the source.
func (t *Tracker) AddBytes(n int64) { t.bytesIn.Add(n) }

// AddLineIn records one line read.
func (t *Tracker) AddLineIn() { t.linesIn.Add(1) }

// AddLineOut records one line written.
func (t *Tracker) AddLineOut() { t.linesOut.Add(1) }

// Summary returns a one-line summary string.
func (t *Tracker) Summary() string {
	return fmt.Sprintf("read %d lines (%.1f KB), wrote %d lines",
		t.linesIn.Load(),
		float64(t.bytesIn.Load())/1024,
		t.linesOut.Load(),
	)
}

// Stop halts the background reporter and prints a final summary.
func (t *Tracker) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		close(t.done)
	}
	if t.out != nil {
		fmt.Fprintln(t.out, "done:", t.Summary())
	}
}

func (t *Tracker) loop() {
	for {
		select {
		case <-t.ticker.C:
			fmt.Fprintln(t.out, "progress:", t.Summary())
		case <-t.done:
			return
		}
	}
}

// Package rotate provides output file rotation based on configurable
// strategies such as maximum file size or maximum line count per output file.
package rotate

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
)

// Strategy controls when a new output file is created.
type Strategy int

const (
	// ByLines rotates after a fixed number of lines.
	ByLines Strategy = iota
	// BySize rotates after a fixed number of bytes.
	BySize
)

// Rotator manages a sequence of output files, opening a new one whenever
// the current file exceeds the configured threshold.
type Rotator struct {
	dir      string
	prefix   string
	strategy Strategy
	threshold int64

	current  *os.File
	fileIdx  int
	counter  int64 // lines written or bytes written to current file
	total    atomic.Int64
}

// New creates a Rotator that writes files under dir with the given prefix.
// threshold is interpreted as a line count (ByLines) or byte count (BySize).
func New(dir, prefix string, strategy Strategy, threshold int64) (*Rotator, error) {
	if threshold <= 0 {
		return nil, fmt.Errorf("rotate: threshold must be > 0, got %d", threshold)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("rotate: mkdir %s: %w", dir, err)
	}
	r := &Rotator{dir: dir, prefix: prefix, strategy: strategy, threshold: threshold}
	if err := r.openNext(); err != nil {
		return nil, err
	}
	return r, nil
}

// WriteLine writes a single log line (without trailing newline) to the current
// file, rotating to a new file if the threshold has been reached.
func (r *Rotator) WriteLine(line string) error {
	data := line + "\n"
	if err := r.rotateIfNeeded(int64(len(data))); err != nil {
		return err
	}
	n, err := fmt.Fprint(r.current, data)
	if err != nil {
		return fmt.Errorf("rotate: write: %w", err)
	}
	r.counter += int64(n)
	r.total.Add(1)
	return nil
}

// Close flushes and closes the current output file.
func (r *Rotator) Close() error {
	if r.current != nil {
		return r.current.Close()
	}
	return nil
}

// FilesCreated returns the total number of output files opened so far.
func (r *Rotator) FilesCreated() int { return r.fileIdx }

// LinesWritten returns the cumulative number of lines written across all files.
func (r *Rotator) LinesWritten() int64 { return r.total.Load() }

func (r *Rotator) rotateIfNeeded(incoming int64) error {
	switch r.strategy {
	case ByLines:
		if r.counter >= r.threshold {
			return r.rotate()
		}
	case BySize:
		if r.counter+incoming > r.threshold {
			return r.rotate()
		}
	}
	return nil
}

func (r *Rotator) rotate() error {
	if err := r.current.Close(); err != nil {
		return fmt.Errorf("rotate: close: %w", err)
	}
	return r.openNext()
}

func (r *Rotator) openNext() error {
	r.fileIdx++
	r.counter = 0
	name := filepath.Join(r.dir, fmt.Sprintf("%s_%04d.log", r.prefix, r.fileIdx))
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("rotate: create %s: %w", name, err)
	}
	r.current = f
	return nil
}

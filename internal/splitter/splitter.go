// Package splitter coordinates reading, filtering, and writing log lines
// within a specified time range from a source reader to a destination writer.
package splitter

import (
	"fmt"
	"io"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/progress"
	"github.com/yourorg/logslice/internal/reader"
	"github.com/yourorg/logslice/internal/writer"
)

// Config holds the parameters for a split operation.
type Config struct {
	Source  io.Reader
	Dest    io.Writer
	Filter  *filter.TimeFilter
	Tracker *progress.Tracker
}

// Splitter reads log lines from a source, applies a time filter, and writes
// matching lines to the destination.
type Splitter struct {
	cfg Config
}

// New creates a new Splitter from the provided Config.
func New(cfg Config) (*Splitter, error) {
	if cfg.Source == nil {
		return nil, fmt.Errorf("splitter: source reader must not be nil")
	}
	if cfg.Dest == nil {
		return nil, fmt.Errorf("splitter: destination writer must not be nil")
	}
	if cfg.Filter == nil {
		return nil, fmt.Errorf("splitter: time filter must not be nil")
	}
	return &Splitter{cfg: cfg}, nil
}

// Run executes the split operation, returning the number of lines written and
// any error encountered.
func (s *Splitter) Run(lr *reader.LineReader, lw *writer.LineWriter) (int, error) {
	written := 0
	for {
		line, err := lr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return written, fmt.Errorf("splitter: read error: %w", err)
		}

		if s.cfg.Tracker != nil {
			s.cfg.Tracker.Inc()
		}

		if !s.cfg.Filter.Match(line) {
			continue
		}

		if err := lw.WriteLine(line); err != nil {
			return written, fmt.Errorf("splitter: write error: %w", err)
		}
		written++
	}
	return written, nil
}

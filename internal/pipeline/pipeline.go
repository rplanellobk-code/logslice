// Package pipeline wires together the reader, filter, and writer components
// into a single cohesive processing pipeline for log extraction.
package pipeline

import (
	"context"
	"fmt"
	"io"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/progress"
	"github.com/logslice/logslice/internal/reader"
	"github.com/logslice/logslice/internal/writer"
)

// Pipeline reads log lines from a source, applies a time filter, and writes
// matching lines to a destination, reporting progress along the way.
type Pipeline struct {
	reader  *reader.LineReader
	filter  *filter.TimeFilter
	writer  *writer.LineWriter
	tracker *progress.Tracker
}

// New constructs a Pipeline from the provided components.
// Returns an error if any required component is nil.
func New(r *reader.LineReader, f *filter.TimeFilter, w *writer.LineWriter, t *progress.Tracker) (*Pipeline, error) {
	if r == nil {
		return nil, fmt.Errorf("pipeline: reader must not be nil")
	}
	if f == nil {
		return nil, fmt.Errorf("pipeline: filter must not be nil")
	}
	if w == nil {
		return nil, fmt.Errorf("pipeline: writer must not be nil")
	}
	return &Pipeline{reader: r, filter: f, writer: w, tracker: t}, nil
}

// Run executes the pipeline until the source is exhausted, the context is
// cancelled, or an unrecoverable error occurs.
// It returns the number of lines written and any error encountered.
func (p *Pipeline) Run(ctx context.Context) (int64, error) {
	var written int64

	for {
		if err := ctx.Err(); err != nil {
			return written, err
		}

		line, err := p.reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return written, fmt.Errorf("pipeline: read error: %w", err)
		}

		if p.tracker != nil {
			p.tracker.Add(1)
		}

		if !p.filter.Match(line) {
			continue
		}

		if err := p.writer.WriteLine(line); err != nil {
			return written, fmt.Errorf("pipeline: write error: %w", err)
		}
		written++
	}

	return written, nil
}

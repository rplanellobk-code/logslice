// Package linenum provides line-number annotation for log lines.
// Each line passing through the annotator is prefixed with a
// monotonically increasing counter so downstream consumers can
// reference exact positions within a processed stream.
package linenum

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/aurc/logslice/internal/reader"
)

// ErrNilLine is returned when a nil LogLine pointer is supplied.
var ErrNilLine = errors.New("linenum: nil log line")

// Annotator prepends a line number to the Raw field of every LogLine
// it processes. The counter starts at 1 and is safe for concurrent use.
type Annotator struct {
	counter atomic.Int64
	format  string
}

// Option is a functional option for Annotator.
type Option func(*Annotator)

// WithFormat overrides the fmt format string used to build the prefix.
// The format must contain exactly one integer verb, e.g. "%06d | ".
func WithFormat(f string) Option {
	return func(a *Annotator) {
		if f != "" {
			a.format = f
		}
	}
}

// New creates an Annotator. The counter starts at zero; the first call
// to Annotate will emit line number 1.
func New(opts ...Option) *Annotator {
	a := &Annotator{format: "%d | "}
	for _, o := range opts {
		o(a)
	}
	return a
}

// Annotate prepends the current line number to line.Raw and increments
// the internal counter. It returns ErrNilLine if line is nil.
func (a *Annotator) Annotate(line *reader.LogLine) error {
	if line == nil {
		return ErrNilLine
	}
	n := a.counter.Add(1)
	line.Raw = fmt.Sprintf(a.format, n) + line.Raw
	return nil
}

// Count returns the number of lines annotated so far.
func (a *Annotator) Count() int64 {
	return a.counter.Load()
}

// Reset sets the internal counter back to zero.
func (a *Annotator) Reset() {
	a.counter.Store(0)
}

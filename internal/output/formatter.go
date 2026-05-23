// Package output handles formatting and writing of extracted log lines
// to one or more destination writers.
package output

import (
	"fmt"
	"io"
	"time"

	"github.com/user/logslice/internal/reader"
)

// Format controls how output lines are rendered.
type Format int

const (
	// FormatRaw writes the original log line unchanged.
	FormatRaw Format = iota
	// FormatNormalized rewrites the timestamp in RFC3339 form.
	FormatNormalized
)

// Formatter writes LogLine values to a destination writer,
// optionally transforming the timestamp representation.
type Formatter struct {
	w      io.Writer
	fmt    Format
	timeFmt string
}

// NewFormatter creates a Formatter that writes to w using the given Format.
// timeFmt is only used when format is FormatNormalized; pass an empty string
// to default to time.RFC3339Nano.
func NewFormatter(w io.Writer, format Format, timeFmt string) (*Formatter, error) {
	if w == nil {
		return nil, fmt.Errorf("output: writer must not be nil")
	}
	if timeFmt == "" {
		timeFmt = time.RFC3339Nano
	}
	return &Formatter{w: w, fmt: format, timeFmt: timeFmt}, nil
}

// Write renders line to the underlying writer.
// It returns the number of bytes written and any error.
func (f *Formatter) Write(line reader.LogLine) (int, error) {
	var text string
	switch f.fmt {
	case FormatNormalized:
		ts := line.Time.Format(f.timeFmt)
		// Replace the original timestamp prefix with the normalised one.
		// We assume the raw line begins with the original timestamp token.
		text = ts + "\t" + line.Raw + "\n"
	default:
		text = line.Raw + "\n"
	}
	return io.WriteString(f.w, text)
}

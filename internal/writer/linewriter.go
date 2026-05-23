package writer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/yourorg/logslice/internal/reader"
)

// LineWriter writes filtered log lines to an output destination.
type LineWriter struct {
	w       *bufio.Writer
	written int
}

// NewLineWriter creates a new LineWriter wrapping the given io.Writer.
func NewLineWriter(w io.Writer) *LineWriter {
	return &LineWriter{
		w: bufio.NewWriter(w),
	}
}

// WriteLine writes a single LogLine's raw content followed by a newline.
func (lw *LineWriter) WriteLine(line reader.LogLine) error {
	_, err := fmt.Fprintln(lw.w, line.Raw)
	if err != nil {
		return fmt.Errorf("linewriter: write failed: %w", err)
	}
	lw.written++
	return nil
}

// WriteLines writes all provided LogLines to the output.
func (lw *LineWriter) WriteLines(lines []reader.LogLine) error {
	for _, line := range lines {
		if err := lw.WriteLine(line); err != nil {
			return err
		}
	}
	return nil
}

// Flush flushes any buffered data to the underlying writer.
func (lw *LineWriter) Flush() error {
	if err := lw.w.Flush(); err != nil {
		return fmt.Errorf("linewriter: flush failed: %w", err)
	}
	return nil
}

// Written returns the number of lines successfully written.
func (lw *LineWriter) Written() int {
	return lw.written
}

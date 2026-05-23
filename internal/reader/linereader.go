package reader

import (
	"bufio"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/timestamp"
)

// LogLine represents a single parsed log line with its extracted timestamp.
type LogLine struct {
	Raw       string
	Timestamp time.Time
	Valid     bool
}

// LineReader reads log lines from an io.Reader and extracts timestamps.
type LineReader struct {
	scanner *bufio.Scanner
	parser  *timestamp.Parser
}

// NewLineReader creates a LineReader that reads from r.
// It uses the provided timestamp.Parser to extract timestamps from each line.
func NewLineReader(r io.Reader, p *timestamp.Parser) *LineReader {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	return &LineReader{
		scanner: scanner,
		parser:  p,
	}
}

// Next advances to the next line and returns a LogLine.
// Returns (LogLine{}, false) when there are no more lines or an error occurs.
func (lr *LineReader) Next() (LogLine, bool) {
	if !lr.scanner.Scan() {
		return LogLine{}, false
	}
	raw := lr.scanner.Text()
	if raw == "" {
		return LogLine{Raw: raw, Valid: false}, true
	}
	t, err := lr.parser.Parse(raw)
	if err != nil {
		return LogLine{Raw: raw, Valid: false}, true
	}
	return LogLine{Raw: raw, Timestamp: t, Valid: true}, true
}

// Err returns any error encountered by the underlying scanner.
func (lr *LineReader) Err() error {
	return lr.scanner.Err()
}

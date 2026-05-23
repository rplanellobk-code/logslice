package reader

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/timestamp"
)

func newTestParser(t *testing.T) *timestamp.Parser {
	t.Helper()
	p, err := timestamp.NewParser("", time.UTC)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	return p
}

func TestLineReaderBasic(t *testing.T) {
	input := "2024-01-15T10:00:00Z INFO starting service\n2024-01-15T10:00:01Z DEBUG connected\n"
	p := newTestParser(t)
	lr := NewLineReader(strings.NewReader(input), p)

	lines := []LogLine{}
	for {
		line, ok := lr.Next()
		if !ok {
			break
		}
		lines = append(lines, line)
	}
	if err := lr.Err(); err != nil {
		t.Fatalf("unexpected scanner error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for _, l := range lines {
		if !l.Valid {
			t.Errorf("expected valid timestamp for line: %q", l.Raw)
		}
	}
}

func TestLineReaderInvalidTimestamp(t *testing.T) {
	input := "no timestamp here\n2024-01-15T10:00:00Z INFO ok\n"
	p := newTestParser(t)
	lr := NewLineReader(strings.NewReader(input), p)

	lines := []LogLine{}
	for {
		line, ok := lr.Next()
		if !ok {
			break
		}
		lines = append(lines, line)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0].Valid {
		t.Errorf("expected first line to be invalid")
	}
	if !lines[1].Valid {
		t.Errorf("expected second line to be valid")
	}
}

func TestLineReaderEmpty(t *testing.T) {
	p := newTestParser(t)
	lr := NewLineReader(strings.NewReader(""), p)
	_, ok := lr.Next()
	if ok {
		t.Error("expected no lines from empty reader")
	}
}

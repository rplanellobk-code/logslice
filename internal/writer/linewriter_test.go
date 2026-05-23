package writer

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/reader"
)

func makeLogLine(raw string) reader.LogLine {
	return reader.LogLine{
		Raw:       raw,
		Timestamp: time.Now(),
	}
}

func TestLineWriterWriteLine(t *testing.T) {
	var buf bytes.Buffer
	lw := NewLineWriter(&buf)

	line := makeLogLine("2024-01-01T00:00:00Z INFO hello world")
	if err := lw.WriteLine(line); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := lw.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "hello world") {
		t.Errorf("expected output to contain 'hello world', got: %q", got)
	}
	if lw.Written() != 1 {
		t.Errorf("expected Written()=1, got %d", lw.Written())
	}
}

func TestLineWriterWriteLines(t *testing.T) {
	var buf bytes.Buffer
	lw := NewLineWriter(&buf)

	lines := []reader.LogLine{
		makeLogLine("line one"),
		makeLogLine("line two"),
		makeLogLine("line three"),
	}

	if err := lw.WriteLines(lines); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := lw.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	output := buf.String()
	for _, l := range lines {
		if !strings.Contains(output, l.Raw) {
			t.Errorf("expected output to contain %q", l.Raw)
		}
	}
	if lw.Written() != 3 {
		t.Errorf("expected Written()=3, got %d", lw.Written())
	}
}

func TestLineWriterEmpty(t *testing.T) {
	var buf bytes.Buffer
	lw := NewLineWriter(&buf)

	if err := lw.WriteLines(nil); err != nil {
		t.Fatalf("unexpected error writing nil slice: %v", err)
	}
	if err := lw.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}
	if lw.Written() != 0 {
		t.Errorf("expected Written()=0, got %d", lw.Written())
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

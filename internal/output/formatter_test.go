package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/reader"
)

func makeLine(raw string, ts time.Time) reader.LogLine {
	return reader.LogLine{Raw: raw, Time: ts}
}

func TestFormatterNilWriter(t *testing.T) {
	_, err := output.NewFormatter(nil, output.FormatRaw, "")
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestFormatterRaw(t *testing.T) {
	var buf bytes.Buffer
	f, err := output.NewFormatter(&buf, output.FormatRaw, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ts := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	line := makeLine("2024-01-15T12:00:00Z INFO hello world", ts)

	n, err := f.Write(line)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if n == 0 {
		t.Fatal("expected non-zero bytes written")
	}

	got := buf.String()
	if !strings.Contains(got, line.Raw) {
		t.Errorf("raw output missing original line; got %q", got)
	}
}

func TestFormatterNormalized(t *testing.T) {
	var buf bytes.Buffer
	f, err := output.NewFormatter(&buf, output.FormatNormalized, time.RFC3339)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ts := time.Date(2024, 6, 1, 8, 30, 0, 0, time.UTC)
	line := makeLine("01/Jun/2024:08:30:00 +0000 GET /api", ts)

	_, err = f.Write(line)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	got := buf.String()
	if !strings.HasPrefix(got, "2024-06-01T08:30:00Z") {
		t.Errorf("normalized output should start with RFC3339 timestamp; got %q", got)
	}
}

func TestFormatterDefaultTimeFmt(t *testing.T) {
	var buf bytes.Buffer
	// Empty timeFmt should default to RFC3339Nano without error.
	f, err := output.NewFormatter(&buf, output.FormatNormalized, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ts := time.Date(2024, 3, 10, 0, 0, 0, 500, time.UTC)
	_, err = f.Write(makeLine("some raw line", ts))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

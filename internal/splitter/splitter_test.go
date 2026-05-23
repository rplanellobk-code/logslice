package splitter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/reader"
	"github.com/yourorg/logslice/internal/splitter"
	"github.com/yourorg/logslice/internal/timestamp"
	"github.com/yourorg/logslice/internal/writer"
)

const logData = `2024-01-10T08:00:00Z INFO startup
2024-01-10T09:00:00Z INFO request received
2024-01-10T10:00:00Z WARN slow query
2024-01-10T11:00:00Z ERROR disk full
2024-01-10T12:00:00Z INFO shutdown
`

func makeComponents(t *testing.T, src string, start, end time.Time) (
	*splitter.Splitter, *reader.LineReader, *writer.LineWriter, *bytes.Buffer,
) {
	t.Helper()
	dst := &bytes.Buffer{}
	parser, err := timestamp.NewParser("", time.UTC)
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}
	lr := reader.NewLineReader(strings.NewReader(src), parser)
	lw := writer.NewLineWriter(dst)
	f, err := filter.NewTimeFilter(start, end)
	if err != nil {
		t.Fatalf("NewTimeFilter: %v", err)
	}
	sp, err := splitter.New(splitter.Config{
		Source: strings.NewReader(src),
		Dest:   dst,
		Filter: f,
	})
	if err != nil {
		t.Fatalf("splitter.New: %v", err)
	}
	return sp, lr, lw, dst
}

func TestSplitterMatchesRange(t *testing.T) {
	start := time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 10, 11, 0, 0, 0, time.UTC)
	sp, lr, lw, dst := makeComponents(t, logData, start, end)

	n, err := sp.Run(lr, lw)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 lines written, got %d", n)
	}
	if !strings.Contains(dst.String(), "request received") {
		t.Error("expected 'request received' in output")
	}
}

func TestSplitterNilSource(t *testing.T) {
	f, _ := filter.NewTimeFilter(time.Time{}, time.Time{})
	_, err := splitter.New(splitter.Config{Source: nil, Dest: &bytes.Buffer{}, Filter: f})
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestSplitterNilDest(t *testing.T) {
	f, _ := filter.NewTimeFilter(time.Time{}, time.Time{})
	_, err := splitter.New(splitter.Config{Source: strings.NewReader(""), Dest: nil, Filter: f})
	if err == nil {
		t.Error("expected error for nil dest")
	}
}

func TestSplitterEmptyInput(t *testing.T) {
	start := time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 10, 11, 0, 0, 0, time.UTC)
	sp, lr, lw, _ := makeComponents(t, "", start, end)
	n, err := sp.Run(lr, lw)
	if err != nil {
		t.Fatalf("Run on empty input: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
}

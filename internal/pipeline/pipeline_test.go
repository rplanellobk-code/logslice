package pipeline_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/pipeline"
	"github.com/logslice/logslice/internal/reader"
	"github.com/logslice/logslice/internal/timestamp"
	"github.com/logslice/logslice/internal/writer"
)

func makeComponents(t *testing.T, input string, from, to time.Time) (*reader.LineReader, *filter.TimeFilter, *writer.LineWriter) {
	t.Helper()
	p, err := timestamp.NewParser("", time.UTC)
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}
	r := reader.NewLineReader(strings.NewReader(input), p)
	f, err := filter.NewTimeFilter(from, to)
	if err != nil {
		t.Fatalf("NewTimeFilter: %v", err)
	}
	var buf bytes.Buffer
	w := writer.NewLineWriter(&buf)
	return r, f, w
}

func TestPipelineRun(t *testing.T) {
	input := "2024-01-01T10:00:00Z line one\n2024-01-01T11:00:00Z line two\n2024-01-01T12:00:00Z line three\n"
	from := time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 11, 30, 0, 0, time.UTC)

	r, f, w := makeComponents(t, input, from, to)
	pl, err := pipeline.New(r, f, w, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	n, err := pl.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 line written, got %d", n)
	}
}

func TestPipelineNilReader(t *testing.T) {
	_, f, w := makeComponents(t, "", time.Time{}, time.Time{})
	_, err := pipeline.New(nil, f, w, nil)
	if err == nil {
		t.Fatal("expected error for nil reader")
	}
}

func TestPipelineNilFilter(t *testing.T) {
	r, _, w := makeComponents(t, "", time.Time{}, time.Time{})
	_, err := pipeline.New(r, nil, w, nil)
	if err == nil {
		t.Fatal("expected error for nil filter")
	}
}

func TestPipelineNilWriter(t *testing.T) {
	r, f, _ := makeComponents(t, "", time.Time{}, time.Time{})
	_, err := pipeline.New(r, f, nil, nil)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestPipelineCancelledContext(t *testing.T) {
	input := "2024-01-01T10:00:00Z line one\n"
	r, f, w := makeComponents(t, input, time.Time{}, time.Time{})
	pl, _ := pipeline.New(r, f, w, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := pl.Run(ctx)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

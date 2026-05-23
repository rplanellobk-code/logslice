package index_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/index"
	"github.com/user/logslice/internal/timestamp"
)

const sampleLogs = `2024-01-01T10:00:00Z INFO starting server
2024-01-01T10:01:00Z INFO listening on :8080
2024-01-01T10:02:00Z WARN high memory usage
2024-01-01T10:03:00Z ERROR disk full
`

func newParser(t *testing.T) *timestamp.Parser {
	t.Helper()
	p, err := timestamp.NewParser("", time.UTC)
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}
	return p
}

func TestBuildIndex(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	p := newParser(t)

	idx, err := index.Build(r, p)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if len(idx.Entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(idx.Entries))
	}
	if idx.Entries[0].Offset != 0 {
		t.Errorf("first entry offset should be 0, got %d", idx.Entries[0].Offset)
	}
}

func TestFindStart(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	p := newParser(t)
	idx, _ := index.Build(r, p)

	start := time.Date(2024, 1, 1, 10, 1, 0, 0, time.UTC)
	offset := idx.FindStart(start)
	if offset != idx.Entries[1].Offset {
		t.Errorf("FindStart: expected offset %d, got %d", idx.Entries[1].Offset, offset)
	}
}

func TestFindStartZero(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	p := newParser(t)
	idx, _ := index.Build(r, p)

	if got := idx.FindStart(time.Time{}); got != 0 {
		t.Errorf("zero start should return 0, got %d", got)
	}
}

func TestFindEnd(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	p := newParser(t)
	idx, _ := index.Build(r, p)

	end := time.Date(2024, 1, 1, 10, 2, 0, 0, time.UTC)
	offset := idx.FindEnd(end)
	if offset != idx.Entries[2].Offset {
		t.Errorf("FindEnd: expected offset %d, got %d", idx.Entries[2].Offset, offset)
	}
}

func TestFindEndZero(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	p := newParser(t)
	idx, _ := index.Build(r, p)

	if got := idx.FindEnd(time.Time{}); got != -1 {
		t.Errorf("zero end should return -1, got %d", got)
	}
}

func TestBuildEmptyReader(t *testing.T) {
	r := strings.NewReader("")
	p := newParser(t)
	idx, err := index.Build(r, p)
	if err != nil {
		t.Fatalf("Build on empty input: %v", err)
	}
	if len(idx.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(idx.Entries))
	}
}

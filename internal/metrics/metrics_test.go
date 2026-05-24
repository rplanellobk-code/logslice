package metrics

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewCountersZero(t *testing.T) {
	c := New()
	if c.LinesRead.Load() != 0 {
		t.Fatalf("expected 0 lines read, got %d", c.LinesRead.Load())
	}
	if c.BytesWritten.Load() != 0 {
		t.Fatalf("expected 0 bytes written, got %d", c.BytesWritten.Load())
	}
}

func TestCountersAtomic(t *testing.T) {
	c := New()
	c.LinesRead.Add(10)
	c.LinesMatched.Add(7)
	c.LinesWritten.Add(7)
	c.LinesSkipped.Add(3)
	c.BytesWritten.Add(1024)

	if c.LinesRead.Load() != 10 {
		t.Errorf("LinesRead: want 10, got %d", c.LinesRead.Load())
	}
	if c.LinesSkipped.Load() != 3 {
		t.Errorf("LinesSkipped: want 3, got %d", c.LinesSkipped.Load())
	}
	if c.BytesWritten.Load() != 1024 {
		t.Errorf("BytesWritten: want 1024, got %d", c.BytesWritten.Load())
	}
}

func TestSummaryContainsFields(t *testing.T) {
	c := New()
	c.LinesRead.Add(100)
	c.LinesMatched.Add(42)
	c.LinesWritten.Add(42)
	c.BytesWritten.Add(8192)

	var buf bytes.Buffer
	c.Summary(&buf)
	out := buf.String()

	for _, want := range []string{"elapsed", "lines read", "lines matched", "bytes written"} {
		if !strings.Contains(out, want) {
			t.Errorf("Summary missing field %q", want)
		}
	}
}

func TestElapsedPositive(t *testing.T) {
	c := New()
	time.Sleep(2 * time.Millisecond)
	if c.Elapsed() <= 0 {
		t.Error("expected positive elapsed duration")
	}
}

func TestRateZeroElapsed(t *testing.T) {
	// Construct a Counters with a future start to force near-zero elapsed.
	c := &Counters{start: time.Now().Add(10 * time.Hour)}
	if c.Rate() != 0 {
		t.Errorf("expected rate 0 for future start, got %f", c.Rate())
	}
}

func TestRatePositive(t *testing.T) {
	c := New()
	c.LinesRead.Add(500)
	time.Sleep(5 * time.Millisecond)
	if c.Rate() <= 0 {
		t.Error("expected positive rate")
	}
}

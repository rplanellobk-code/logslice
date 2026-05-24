package transform_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/reader"
	"github.com/yourorg/logslice/internal/transform"
)

func makeLine(text string) reader.LogLine {
	return reader.LogLine{Timestamp: time.Unix(1_000_000, 0), Text: text}
}

func TestNewNilFuncReturnsError(t *testing.T) {
	_, err := transform.New(nil)
	if err == nil {
		t.Fatal("expected error for nil Func, got nil")
	}
}

func TestNewEmptyChainOK(t *testing.T) {
	tr, err := transform.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeLine("hello")
	got := tr.Apply(line)
	if got.Text != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got.Text)
	}
}

func TestApplySingleFunc(t *testing.T) {
	tr, err := transform.New(strings.ToUpper)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := tr.Apply(makeLine("hello world"))
	if got.Text != "HELLO WORLD" {
		t.Fatalf("expected %q, got %q", "HELLO WORLD", got.Text)
	}
}

func TestApplyChainOrder(t *testing.T) {
	prefix := func(s string) string { return "[A]" + s }
	suffix := func(s string) string { return s + "[B]" }
	tr, _ := transform.New(prefix, suffix)
	got := tr.Apply(makeLine("x"))
	want := "[A]x[B]"
	if got.Text != want {
		t.Fatalf("expected %q, got %q", want, got.Text)
	}
}

func TestApplyPreservesTimestamp(t *testing.T) {
	ts := time.Unix(42, 0)
	line := reader.LogLine{Timestamp: ts, Text: "msg"}
	tr, _ := transform.New(strings.ToUpper)
	got := tr.Apply(line)
	if !got.Timestamp.Equal(ts) {
		t.Fatalf("timestamp changed: got %v, want %v", got.Timestamp, ts)
	}
}

func TestApplyAllLength(t *testing.T) {
	tr, _ := transform.New(strings.TrimSpace)
	lines := []reader.LogLine{makeLine("  a  "), makeLine(" b "), makeLine("c")}
	out := tr.ApplyAll(lines)
	if len(out) != len(lines) {
		t.Fatalf("expected %d lines, got %d", len(lines), len(out))
	}
	for i, l := range out {
		if strings.Contains(l.Text, " ") {
			t.Errorf("line %d still has spaces: %q", i, l.Text)
		}
	}
}

func TestApplyAllEmptySlice(t *testing.T) {
	tr, _ := transform.New(strings.ToUpper)
	out := tr.ApplyAll(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d elements", len(out))
	}
}

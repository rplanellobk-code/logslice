package prefix_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/prefix"
	"github.com/logslice/logslice/internal/reader"
)

func makeLine(raw string) reader.LogLine {
	return reader.LogLine{Raw: raw, Timestamp: time.Now()}
}

func TestNewEmptyPrefix(t *testing.T) {
	_, err := prefix.New("")
	if err == nil {
		t.Fatal("expected error for empty prefix, got nil")
	}
}

func TestNewWhitespaceOnlyPrefix(t *testing.T) {
	_, err := prefix.New("   ")
	if err == nil {
		t.Fatal("expected error for whitespace-only prefix, got nil")
	}
}

func TestApplyDefaultSeparator(t *testing.T) {
	p, err := prefix.New("[INFO]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := p.Apply(makeLine("hello world"))
	want := "[INFO] hello world"
	if got.Raw != want {
		t.Errorf("Apply() = %q, want %q", got.Raw, want)
	}
}

func TestApplyCustomSeparator(t *testing.T) {
	p, err := prefix.New("SVC", prefix.WithSeparator(" | "))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := p.Apply(makeLine("request received"))
	want := "SVC | request received"
	if got.Raw != want {
		t.Errorf("Apply() = %q, want %q", got.Raw, want)
	}
}

func TestApplyPreservesTimestamp(t *testing.T) {
	p, _ := prefix.New("TAG")
	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	line := reader.LogLine{Raw: "msg", Timestamp: ts}
	got := p.Apply(line)
	if !got.Timestamp.Equal(ts) {
		t.Errorf("Timestamp changed: got %v, want %v", got.Timestamp, ts)
	}
}

func TestApplyEmptyLine(t *testing.T) {
	p, _ := prefix.New("PRE")
	got := p.Apply(makeLine(""))
	want := "PRE "
	if got.Raw != want {
		t.Errorf("Apply() = %q, want %q", got.Raw, want)
	}
}

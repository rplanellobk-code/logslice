package head

import (
	"strings"
	"testing"
)

func TestNewInvalidN(t *testing.T) {
	_, err := New(strings.NewReader(""), 0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
	_, err = New(strings.NewReader(""), -5)
	if err == nil {
		t.Fatal("expected error for n=-5")
	}
}

func TestNewNilReader(t *testing.T) {
	_, err := New(nil, 3)
	if err == nil {
		t.Fatal("expected error for nil reader")
	}
}

func TestNewValidN(t *testing.T) {
	h, err := New(strings.NewReader(""), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil Reader")
	}
}

func TestReadLineReturnsFirstN(t *testing.T) {
	input := "alpha\nbeta\ngamma\ndelta\n"
	h, _ := New(strings.NewReader(input), 2)

	line, ok := h.ReadLine()
	if !ok || line != "alpha" {
		t.Fatalf("expected (alpha, true), got (%q, %v)", line, ok)
	}
	line, ok = h.ReadLine()
	if !ok || line != "beta" {
		t.Fatalf("expected (beta, true), got (%q, %v)", line, ok)
	}
	_, ok = h.ReadLine()
	if ok {
		t.Fatal("expected false after limit reached")
	}
}

func TestLinesCollectsAll(t *testing.T) {
	input := "one\ntwo\nthree\n"
	h, _ := New(strings.NewReader(input), 10)
	lines := h.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "one" || lines[2] != "three" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

func TestLinesRespectsLimit(t *testing.T) {
	input := "a\nb\nc\nd\ne\n"
	h, _ := New(strings.NewReader(input), 3)
	lines := h.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestSeenCounter(t *testing.T) {
	input := "x\ny\n"
	h, _ := New(strings.NewReader(input), 5)
	h.ReadLine()
	if h.Seen() != 1 {
		t.Fatalf("expected Seen()=1, got %d", h.Seen())
	}
	h.ReadLine()
	if h.Seen() != 2 {
		t.Fatalf("expected Seen()=2, got %d", h.Seen())
	}
}

func TestEmptyInput(t *testing.T) {
	h, _ := New(strings.NewReader(""), 5)
	lines := h.Lines()
	if len(lines) != 0 {
		t.Fatalf("expected 0 lines from empty input, got %d", len(lines))
	}
}

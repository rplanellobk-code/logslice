package tail

import (
	"strings"
	"testing"
)

func TestNewInvalidN(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
	_, err = New(-5)
	if err == nil {
		t.Fatal("expected error for n=-5")
	}
}

func TestNewValidN(t *testing.T) {
	tl, err := New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tl == nil {
		t.Fatal("expected non-nil tailer")
	}
}

func TestReadNilReader(t *testing.T) {
	tl, _ := New(5)
	_, err := tl.Read(nil)
	if err == nil {
		t.Fatal("expected error for nil reader")
	}
}

func TestReadEmpty(t *testing.T) {
	tl, _ := New(5)
	lines, err := tl.Read(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
}

func TestReadFewerLinesThanN(t *testing.T) {
	tl, _ := New(10)
	input := "alpha\nbeta\ngamma\n"
	lines, err := tl.Read(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "alpha" || lines[2] != "gamma" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

func TestReadExactlyN(t *testing.T) {
	tl, _ := New(3)
	input := "a\nb\nc\n"
	lines, err := tl.Read(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3, got %d", len(lines))
	}
}

func TestReadMoreLinesThanN(t *testing.T) {
	tl, _ := New(3)
	input := "line1\nline2\nline3\nline4\nline5\n"
	lines, err := tl.Read(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "line3" || lines[1] != "line4" || lines[2] != "line5" {
		t.Fatalf("unexpected tail lines: %v", lines)
	}
}

func TestReadNoTrailingNewline(t *testing.T) {
	tl, _ := New(2)
	input := "x\ny\nz"
	lines, err := tl.Read(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "y" || lines[1] != "z" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

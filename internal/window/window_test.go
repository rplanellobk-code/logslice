package window

import (
	"testing"
)

func TestNewInvalidCapacity(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for capacity 0")
	}
	_, err = New(-3)
	if err == nil {
		t.Fatal("expected error for negative capacity")
	}
}

func TestNewValidCapacity(t *testing.T) {
	b, err := New(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Len() != 0 {
		t.Fatalf("expected empty buffer, got len %d", b.Len())
	}
}

func TestPushBelowCapacity(t *testing.T) {
	b, _ := New(5)
	b.Push("a")
	b.Push("b")
	lines := b.Lines()
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "a" || lines[1] != "b" {
		t.Fatalf("unexpected order: %v", lines)
	}
}

func TestPushEvictsOldest(t *testing.T) {
	b, _ := New(3)
	for _, s := range []string{"a", "b", "c", "d"} {
		b.Push(s)
	}
	lines := b.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "b" || lines[1] != "c" || lines[2] != "d" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

func TestLinesEmptyBuffer(t *testing.T) {
	b, _ := New(4)
	if b.Lines() != nil {
		t.Fatal("expected nil from empty buffer")
	}
}

func TestReset(t *testing.T) {
	b, _ := New(3)
	b.Push("x")
	b.Push("y")
	b.Reset()
	if b.Len() != 0 {
		t.Fatalf("expected len 0 after reset, got %d", b.Len())
	}
	if b.Lines() != nil {
		t.Fatal("expected nil lines after reset")
	}
}

func TestWrapAroundMultipleTimes(t *testing.T) {
	b, _ := New(2)
	for i := 0; i < 10; i++ {
		b.Push(string(rune('a' + i)))
	}
	// last two pushes: 'i' (index 8) and 'j' (index 9)
	lines := b.Lines()
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "i" || lines[1] != "j" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

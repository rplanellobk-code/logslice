package dedupe

import (
	"testing"
)

func TestNewInvalidCapacity(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for capacity=0, got nil")
	}
	_, err = New(-5)
	if err == nil {
		t.Fatal("expected error for negative capacity, got nil")
	}
}

func TestIsDuplicateFreshLines(t *testing.T) {
	f, err := New(10)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	lines := []string{"alpha", "beta", "gamma"}
	for _, l := range lines {
		if f.IsDuplicate(l) {
			t.Errorf("line %q reported as duplicate on first encounter", l)
		}
	}
	if f.Len() != 3 {
		t.Errorf("expected Len 3, got %d", f.Len())
	}
}

func TestIsDuplicateDetected(t *testing.T) {
	f, _ := New(10)
	f.IsDuplicate("hello")
	if !f.IsDuplicate("hello") {
		t.Error("expected duplicate to be detected")
	}
}

func TestWindowEviction(t *testing.T) {
	f, _ := New(3)
	for _, l := range []string{"a", "b", "c"} {
		f.IsDuplicate(l)
	}
	// Window is full: adding "d" should evict "a".
	f.IsDuplicate("d")
	if f.Len() != 3 {
		t.Errorf("expected window size 3 after eviction, got %d", f.Len())
	}
	// "a" should no longer be considered a duplicate.
	if f.IsDuplicate("a") {
		t.Error("evicted line 'a' should not be a duplicate")
	}
}

func TestReset(t *testing.T) {
	f, _ := New(5)
	f.IsDuplicate("x")
	f.IsDuplicate("y")
	f.Reset()
	if f.Len() != 0 {
		t.Errorf("expected Len 0 after Reset, got %d", f.Len())
	}
	if f.IsDuplicate("x") {
		t.Error("line 'x' should not be duplicate after Reset")
	}
}

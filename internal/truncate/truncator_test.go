package truncate

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewInvalidMaxLen(t *testing.T) {
	_, err := New(0, false)
	if err == nil {
		t.Fatal("expected error for maxLen=0")
	}
	_, err = New(-5, false)
	if err == nil {
		t.Fatal("expected error for negative maxLen")
	}
}

func TestNewSuffixTooLong(t *testing.T) {
	// DefaultSuffix is 3 bytes; maxLen must be > 3
	_, err := New(3, true)
	if err == nil {
		t.Fatal("expected error when maxLen <= suffix length")
	}
	_, err = New(4, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplyShortLine(t *testing.T) {
	tr, err := New(80, false)
	if err != nil {
		t.Fatal(err)
	}
	line := []byte("short line")
	out := tr.Apply(line)
	if !bytes.Equal(out, line) {
		t.Fatalf("expected unchanged line, got %q", out)
	}
}

func TestApplyExactLength(t *testing.T) {
	tr, _ := New(10, false)
	line := []byte("1234567890")
	out := tr.Apply(line)
	if !bytes.Equal(out, line) {
		t.Fatalf("expected unchanged line, got %q", out)
	}
}

func TestApplyTruncatesWithoutSuffix(t *testing.T) {
	tr, _ := New(10, false)
	line := []byte("1234567890EXTRA")
	out := tr.Apply(line)
	if len(out) != 10 {
		t.Fatalf("expected length 10, got %d", len(out))
	}
	if string(out) != "1234567890" {
		t.Fatalf("unexpected content: %q", out)
	}
}

func TestApplyTruncatesWithSuffix(t *testing.T) {
	tr, _ := New(10, true)
	line := []byte("1234567890EXTRA")
	out := tr.Apply(line)
	if len(out) != 10 {
		t.Fatalf("expected length 10, got %d", len(out))
	}
	if !strings.HasSuffix(string(out), DefaultSuffix) {
		t.Fatalf("expected suffix %q in output %q", DefaultSuffix, out)
	}
}

func TestMaxLen(t *testing.T) {
	tr, _ := New(42, false)
	if tr.MaxLen() != 42 {
		t.Fatalf("expected MaxLen 42, got %d", tr.MaxLen())
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	tr, _ := New(5, true)
	original := []byte("hello world")
	input := make([]byte, len(original))
	copy(input, original)
	tr.Apply(input)
	if !bytes.Equal(input, original) {
		t.Fatal("Apply must not mutate the input slice")
	}
}

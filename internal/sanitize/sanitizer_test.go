package sanitize

import (
	"testing"
)

func TestApplyPrintableUnchanged(t *testing.T) {
	s := New()
	input := "2024-01-02T15:04:05Z INFO hello world"
	if got := s.Apply(input); got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestApplyTabPreserved(t *testing.T) {
	s := New()
	input := "col1\tcol2\tcol3"
	if got := s.Apply(input); got != input {
		t.Fatalf("tab should be preserved, got %q", got)
	}
}

func TestApplyReplacesControlChars(t *testing.T) {
	s := New(WithReplacement('?'))
	input := "hello\x01world\x00"
	want := "hello?world?"
	if got := s.Apply(input); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApplyDefaultReplacement(t *testing.T) {
	s := New()
	input := "bad\x07byte"
	want := "bad\uFFFDbyte"
	if got := s.Apply(input); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestStripANSIDisabledByDefault(t *testing.T) {
	s := New(WithReplacement('?'))
	// ESC is 0x1B — non-printable, so it gets replaced, but the sequence
	// is not stripped as a unit.
	input := "\x1B[31mred\x1B[0m"
	got := s.Apply(input)
	// We just verify the result does not contain the raw ESC byte.
	for _, r := range got {
		if r == 0x1B {
			t.Fatal("ESC byte should have been replaced")
		}
	}
}

func TestStripANSIEnabled(t *testing.T) {
	s := New(WithStripANSI(true))
	input := "\x1B[31mred text\x1B[0m"
	want := "red text"
	if got := s.Apply(input); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestStripANSIMultipleSequences(t *testing.T) {
	s := New(WithStripANSI(true))
	input := "\x1B[1m\x1B[32mgreen bold\x1B[0m normal"
	want := "green bold normal"
	if got := s.Apply(input); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApplyEmptyString(t *testing.T) {
	s := New()
	if got := s.Apply(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

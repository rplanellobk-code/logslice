package highlight

import (
	"strings"
	"testing"
)

func TestNewEmptyPattern(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNewInvalidPattern(t *testing.T) {
	_, err := New("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNewValidPattern(t *testing.T) {
	h, err := New(`ERROR`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil highlighter")
	}
}

func TestApplyNoMatch(t *testing.T) {
	h, _ := New(`ERROR`)
	line := "everything is fine"
	got := h.Apply(line)
	if got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApplyMatchWrapped(t *testing.T) {
	h, _ := New(`ERROR`)
	line := "2024-01-01 ERROR something failed"
	got := h.Apply(line)
	if !strings.Contains(got, Yellow+"ERROR"+Reset) {
		t.Errorf("expected ANSI-wrapped match, got %q", got)
	}
}

func TestApplyCustomColour(t *testing.T) {
	h, _ := New(`WARN`, WithColour(Cyan))
	got := h.Apply("WARN: disk usage high")
	if !strings.Contains(got, Cyan+"WARN"+Reset) {
		t.Errorf("expected cyan highlight, got %q", got)
	}
}

func TestApplyMultipleMatches(t *testing.T) {
	h, _ := New(`\d+`)
	got := h.Apply("retried 3 times after 10 seconds")
	count := strings.Count(got, Reset)
	if count != 2 {
		t.Errorf("expected 2 highlighted matches, got %d", count)
	}
}

func TestApplyPreservesRestOfLine(t *testing.T) {
	h, _ := New(`ERROR`)
	line := "prefix ERROR suffix"
	got := h.Apply(line)
	if !strings.HasPrefix(got, "prefix ") {
		t.Errorf("prefix not preserved: %q", got)
	}
	if !strings.HasSuffix(got, " suffix") {
		t.Errorf("suffix not preserved: %q", got)
	}
}

func TestApplyEmptyLine(t *testing.T) {
	h, _ := New(`ERROR`)
	got := h.Apply("")
	if got != "" {
		t.Errorf("expected empty string for empty input, got %q", got)
	}
}

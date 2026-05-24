package pattern_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/pattern"
)

func TestNewEmptyPattern(t *testing.T) {
	_, err := pattern.New("")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNewInvalidRegex(t *testing.T) {
	_, err := pattern.New("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestMatchFound(t *testing.T) {
	m, err := pattern.New(`ERROR`)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Match("2024-01-01 ERROR something went wrong") {
		t.Error("expected match")
	}
}

func TestMatchNotFound(t *testing.T) {
	m, err := pattern.New(`ERROR`)
	if err != nil {
		t.Fatal(err)
	}
	if m.Match("2024-01-01 INFO all good") {
		t.Error("expected no match")
	}
}

func TestMatchInvert(t *testing.T) {
	m, err := pattern.New(`DEBUG`, pattern.WithInvert())
	if err != nil {
		t.Fatal(err)
	}
	if !m.Match("2024-01-01 INFO hello") {
		t.Error("inverted: non-debug line should pass")
	}
	if m.Match("2024-01-01 DEBUG verbose") {
		t.Error("inverted: debug line should be dropped")
	}
}

func TestMatchRegexCapture(t *testing.T) {
	m, err := pattern.New(`user_id=\d+`)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Match("login user_id=42 success") {
		t.Error("expected match on numeric user_id")
	}
}

func TestPatternReturnsExpression(t *testing.T) {
	expr := `foo.*bar`
	m, _ := pattern.New(expr)
	if m.Pattern() != expr {
		t.Errorf("Pattern() = %q, want %q", m.Pattern(), expr)
	}
}

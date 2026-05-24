package mask

import (
	"strings"
	"testing"
)

func TestNewEmptyPatterns(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}

func TestNewInvalidPattern(t *testing.T) {
	_, err := New([]string{"["})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestApplyDefaultMask(t *testing.T) {
	m, err := New([]string{`\d{4}-\d{4}-\d{4}-\d{4}`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := "card=1234-5678-9012-3456 user=alice"
	got := m.Apply(input)
	if strings.Contains(got, "1234-5678-9012-3456") {
		t.Errorf("sensitive data not redacted: %s", got)
	}
	if !strings.Contains(got, defaultMask) {
		t.Errorf("expected default mask in output: %s", got)
	}
}

func TestApplyCustomMask(t *testing.T) {
	m, err := New([]string{`password=\S+`}, WithMask("[HIDDEN]"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := m.Apply("password=s3cr3t host=localhost")
	if strings.Contains(got, "s3cr3t") {
		t.Errorf("password not masked: %s", got)
	}
	if !strings.Contains(got, "[HIDDEN]") {
		t.Errorf("custom mask not present: %s", got)
	}
}

func TestApplyMultiplePatterns(t *testing.T) {
	m, err := New([]string{
		`\b[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,}\b(?i)`,
		`token=[A-Za-z0-9]+`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := "user=alice@example.com token=abc123"
	got := m.Apply(input)
	if strings.Contains(got, "abc123") {
		t.Errorf("token not masked: %s", got)
	}
}

func TestApplyNoMatch(t *testing.T) {
	m, err := New([]string{`NOMATCH`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := "nothing sensitive here"
	got := m.Apply(input)
	if got != input {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

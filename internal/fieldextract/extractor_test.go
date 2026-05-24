package fieldextract

import (
	"testing"
)

func TestNewInvalidField(t *testing.T) {
	_, err := New("", "=")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewInvalidDelimiter(t *testing.T) {
	_, err := New("level", "")
	if err == nil {
		t.Fatal("expected error for empty delimiter")
	}
}

func TestExtractFound(t *testing.T) {
	ex, _ := New("level", "=")
	val, ok := ex.Extract(`ts=2024-01-01T00:00:00Z level=info msg="hello world"`)
	if !ok {
		t.Fatal("expected field to be found")
	}
	if val != "info" {
		t.Fatalf("expected 'info', got %q", val)
	}
}

func TestExtractNotFound(t *testing.T) {
	ex, _ := New("level", "=")
	_, ok := ex.Extract("ts=2024-01-01T00:00:00Z msg=hello")
	if ok {
		t.Fatal("expected field not to be found")
	}
}

func TestExtractAtEndOfLine(t *testing.T) {
	ex, _ := New("code", "=")
	val, ok := ex.Extract("msg=timeout code=503")
	if !ok {
		t.Fatal("expected field to be found")
	}
	if val != "503" {
		t.Fatalf("expected '503', got %q", val)
	}
}

func TestExtractEmptyValue(t *testing.T) {
	ex, _ := New("key", "=")
	_, ok := ex.Extract("key= other=val")
	if ok {
		t.Fatal("expected empty value to return false")
	}
}

func TestExtractColonDelimiter(t *testing.T) {
	ex, _ := New("host", ":")
	val, ok := ex.Extract("host:web-01 status:200")
	if !ok {
		t.Fatal("expected field to be found")
	}
	if val != "web-01" {
		t.Fatalf("expected 'web-01', got %q", val)
	}
}

func TestFieldAndDelimiterAccessors(t *testing.T) {
	ex, _ := New("svc", "=")
	if ex.Field() != "svc" {
		t.Fatalf("unexpected field %q", ex.Field())
	}
	if ex.Delimiter() != "=" {
		t.Fatalf("unexpected delimiter %q", ex.Delimiter())
	}
}

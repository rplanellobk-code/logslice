package fieldextract

import (
	"testing"
)

func TestNewMultiEmpty(t *testing.T) {
	_, err := NewMulti(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty fields map")
	}
}

func TestNewMultiInvalidField(t *testing.T) {
	_, err := NewMulti(map[string]string{"": "="})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestExtractAllFindsMultiple(t *testing.T) {
	m, err := NewMulti(map[string]string{
		"level": "=",
		"code":  "=",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := m.ExtractAll("ts=2024-01-01 level=warn code=404 msg=notfound")
	if result["level"] != "warn" {
		t.Errorf("expected level=warn, got %q", result["level"])
	}
	if result["code"] != "404" {
		t.Errorf("expected code=404, got %q", result["code"])
	}
}

func TestExtractAllMissingField(t *testing.T) {
	m, _ := NewMulti(map[string]string{
		"level": "=",
		"host":  "=",
	})
	result := m.ExtractAll("level=info msg=ok")
	if _, ok := result["host"]; ok {
		t.Error("host should not be present in result")
	}
	if result["level"] != "info" {
		t.Errorf("expected level=info, got %q", result["level"])
	}
}

func TestExtractAllEmptyLine(t *testing.T) {
	m, _ := NewMulti(map[string]string{"level": "="})
	result := m.ExtractAll("")
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestFieldsReturnsNames(t *testing.T) {
	m, _ := NewMulti(map[string]string{"svc": "=", "env": ":"})
	names := m.Fields()
	if len(names) != 2 {
		t.Fatalf("expected 2 field names, got %d", len(names))
	}
}

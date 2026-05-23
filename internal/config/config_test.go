package config

import (
	"testing"
	"time"
)

func baseConfig() Config {
	return Config{
		Pattern:   "/var/log/*.log",
		OutputDir: "/tmp/out",
	}
}

func TestValidateOK(t *testing.T) {
	cfg := baseConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateEmptyPattern(t *testing.T) {
	cfg := baseConfig()
	cfg.Pattern = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestValidateEmptyOutputDir(t *testing.T) {
	cfg := baseConfig()
	cfg.OutputDir = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty output dir")
	}
}

func TestValidateToBeforeFrom(t *testing.T) {
	cfg := baseConfig()
	now := time.Now()
	cfg.From = now
	cfg.To = now.Add(-time.Hour)
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when 'to' is before 'from'")
	}
}

func TestValidateEqualFromTo(t *testing.T) {
	cfg := baseConfig()
	now := time.Now()
	cfg.From = now
	cfg.To = now
	if err := cfg.Validate(); err != nil {
		t.Fatalf("equal from/to should be valid, got: %v", err)
	}
}

func TestValidateOpenBounds(t *testing.T) {
	cfg := baseConfig()
	// Both bounds zero — open range — should be valid.
	if err := cfg.Validate(); err != nil {
		t.Fatalf("open range should be valid, got: %v", err)
	}
}

func TestValidateOnlyFromSet(t *testing.T) {
	cfg := baseConfig()
	cfg.From = time.Now()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("only 'from' set should be valid, got: %v", err)
	}
}

func TestValidateOnlyToSet(t *testing.T) {
	cfg := baseConfig()
	cfg.To = time.Now()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("only 'to' set should be valid, got: %v", err)
	}
}

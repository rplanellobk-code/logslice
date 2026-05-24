package config

import (
	"testing"
	"time"
)

func baseConfig() Config {
	return Config{
		Pattern:   "/logs/*.log",
		OutputDir: "/out",
	}
}

func TestValidateOK(t *testing.T) {
	if err := Validate(baseConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateEmptyPattern(t *testing.T) {
	c := baseConfig()
	c.Pattern = ""
	if err := Validate(c); err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestValidateEmptyOutputDir(t *testing.T) {
	c := baseConfig()
	c.OutputDir = ""
	if err := Validate(c); err == nil {
		t.Fatal("expected error for empty output_dir")
	}
}

func TestValidateToBeforeFrom(t *testing.T) {
	c := baseConfig()
	c.From = time.Now()
	c.To = c.From.Add(-time.Hour)
	if err := Validate(c); err == nil {
		t.Fatal("expected error when to is before from")
	}
}

func TestValidateToEqualsFrom(t *testing.T) {
	c := baseConfig()
	now := time.Now()
	c.From = now
	c.To = now
	if err := Validate(c); err == nil {
		t.Fatal("expected error when to equals from")
	}
}

func TestValidateOpenRange(t *testing.T) {
	c := baseConfig()
	c.From = time.Now()
	// To is zero — open-ended range, should be valid
	if err := Validate(c); err != nil {
		t.Fatalf("unexpected error for open range: %v", err)
	}
}

func TestValidateNegativeDedupeWindow(t *testing.T) {
	c := baseConfig()
	c.DedupeWindow = -1
	if err := Validate(c); err == nil {
		t.Fatal("expected error for negative dedupe_window")
	}
}

func TestValidateNegativeRateLimit(t *testing.T) {
	c := baseConfig()
	c.RateLimit = -5
	if err := Validate(c); err == nil {
		t.Fatal("expected error for negative rate_limit")
	}
}

package pattern_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/pattern"
)

func TestNewLineFilterNilMatcher(t *testing.T) {
	_, err := pattern.NewLineFilter(nil)
	if err == nil {
		t.Fatal("expected error for nil matcher")
	}
}

func TestLineFilterKeepMatch(t *testing.T) {
	m, _ := pattern.New(`WARN`)
	f, err := pattern.NewLineFilter(m)
	if err != nil {
		t.Fatal(err)
	}
	if !f.Keep("2024-01-01 WARN disk almost full") {
		t.Error("expected Keep=true for matching line")
	}
}

func TestLineFilterDropNonMatch(t *testing.T) {
	m, _ := pattern.New(`WARN`)
	f, _ := pattern.NewLineFilter(m)
	if f.Keep("2024-01-01 INFO all clear") {
		t.Error("expected Keep=false for non-matching line")
	}
}

func TestLineFilterInvertedKeep(t *testing.T) {
	m, _ := pattern.New(`DEBUG`, pattern.WithInvert())
	f, _ := pattern.NewLineFilter(m)
	if !f.Keep("2024-01-01 INFO startup complete") {
		t.Error("inverted filter: INFO line should be kept")
	}
	if f.Keep("2024-01-01 DEBUG trace data") {
		t.Error("inverted filter: DEBUG line should be dropped")
	}
}

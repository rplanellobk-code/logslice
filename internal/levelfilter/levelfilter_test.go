package levelfilter

import (
	"testing"
)

func TestNewUnknownLevel(t *testing.T) {
	_, err := New("verbose")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestNewCaseInsensitive(t *testing.T) {
	f, err := New("WARN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.MinLevel() != LevelWarn {
		t.Fatalf("expected LevelWarn, got %d", f.MinLevel())
	}
}

func TestKeepAtMinLevel(t *testing.T) {
	f, _ := New("warn")
	lines := []struct {
		line string
		want bool
	}{
		{"2024-01-01 DEBUG starting up", false},
		{"2024-01-01 INFO service ready", false},
		{"2024-01-01 WARN disk usage high", true},
		{"2024-01-01 ERROR connection refused", true},
		{"2024-01-01 FATAL out of memory", true},
	}
	for _, tc := range lines {
		got := f.Keep(tc.line)
		if got != tc.want {
			t.Errorf("Keep(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}

func TestKeepDebugPassesAll(t *testing.T) {
	f, _ := New("debug")
	for _, line := range []string{
		"DEBUG init",
		"INFO ready",
		"WARN high load",
		"ERROR crash",
		"FATAL oom",
	} {
		if !f.Keep(line) {
			t.Errorf("Keep(%q) should be true at debug level", line)
		}
	}
}

func TestKeepFatalPassesOnlyFatal(t *testing.T) {
	f, _ := New("fatal")
	if f.Keep("ERROR something bad") {
		t.Error("ERROR should not pass a fatal-level filter")
	}
	if !f.Keep("FATAL system halted") {
		t.Error("FATAL should pass a fatal-level filter")
	}
	if !f.Keep("CRIT kernel panic") {
		t.Error("CRIT alias should pass a fatal-level filter")
	}
}

func TestKeepNoLevelToken(t *testing.T) {
	f, _ := New("info")
	if f.Keep("plain log line with no level token") {
		t.Error("line without level token should not pass")
	}
}

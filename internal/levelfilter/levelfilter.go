// Package levelfilter provides filtering of log lines by severity level.
package levelfilter

import (
	"fmt"
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"warning": LevelWarn,
	"error": LevelError,
	"err":   LevelError,
	"fatal": LevelFatal,
	"crit":  LevelFatal,
}

// Filter discards log lines whose severity is below a minimum level.
type Filter struct {
	minLevel Level
	keywords []string
}

// New creates a Filter that passes only lines at or above minLevelStr.
// minLevelStr is case-insensitive (e.g. "warn", "ERROR").
func New(minLevelStr string) (*Filter, error) {
	norm := strings.ToLower(strings.TrimSpace(minLevelStr))
	lvl, ok := levelNames[norm]
	if !ok {
		return nil, fmt.Errorf("levelfilter: unknown level %q", minLevelStr)
	}
	// Collect all level keywords whose numeric value >= lvl.
	var keywords []string
	seen := map[Level]bool{}
	for name, l := range levelNames {
		if l >= lvl && !seen[l] {
			keywords = append(keywords, name)
			seen[l] = true
		}
	}
	return &Filter{minLevel: lvl, keywords: keywords}, nil
}

// Keep returns true when the line contains a recognised level token that
// is at or above the filter's minimum level.
func (f *Filter) Keep(line string) bool {
	lower := strings.ToLower(line)
	for _, kw := range f.keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// MinLevel returns the configured minimum Level value.
func (f *Filter) MinLevel() Level { return f.minLevel }

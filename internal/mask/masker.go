// Package mask provides field-level redaction for sensitive log data.
// It replaces matched substrings with a configurable mask string.
package mask

import (
	"errors"
	"fmt"
	"regexp"
)

const defaultMask = "***REDACTED***"

// Masker replaces sensitive patterns within a log line.
type Masker struct {
	patterns []*regexp.Regexp
	mask     string
}

// Option configures a Masker.
type Option func(*Masker)

// WithMask overrides the default replacement string.
func WithMask(m string) Option {
	return func(mk *Masker) {
		mk.mask = m
	}
}

// New creates a Masker that replaces any match of the supplied regular
// expression patterns with the mask string.
// Returns an error if patterns is empty or any pattern fails to compile.
func New(patterns []string, opts ...Option) (*Masker, error) {
	if len(patterns) == 0 {
		return nil, errors.New("mask: at least one pattern is required")
	}

	mk := &Masker{mask: defaultMask}
	for _, o := range opts {
		o(mk)
	}

	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("mask: invalid pattern %q: %w", p, err)
		}
		mk.patterns = append(mk.patterns, re)
	}
	return mk, nil
}

// Apply replaces all pattern matches in line with the mask string.
func (m *Masker) Apply(line string) string {
	for _, re := range m.patterns {
		line = re.ReplaceAllString(line, m.mask)
	}
	return line
}

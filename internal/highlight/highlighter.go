package highlight

import (
	"fmt"
	"regexp"
	"strings"
)

// ANSI colour codes used for highlighting.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
)

// Highlighter wraps matched substrings in a log line with ANSI escape codes.
type Highlighter struct {
	pattern *regexp.Regexp
	colour  string
}

// Option is a functional option for Highlighter.
type Option func(*Highlighter)

// WithColour sets the ANSI colour code used for matches.
func WithColour(code string) Option {
	return func(h *Highlighter) {
		h.colour = code
	}
}

// New compiles pattern and returns a Highlighter.
// Returns an error if pattern is empty or invalid.
func New(pattern string, opts ...Option) (*Highlighter, error) {
	if strings.TrimSpace(pattern) == "" {
		return nil, fmt.Errorf("highlight: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("highlight: invalid pattern: %w", err)
	}
	h := &Highlighter{
		pattern: re,
		colour:  Yellow,
	}
	for _, o := range opts {
		o(h)
	}
	return h, nil
}

// Apply returns a copy of line with all pattern matches wrapped in the
// configured ANSI colour codes. If there are no matches the original
// line is returned unchanged.
func (h *Highlighter) Apply(line string) string {
	return h.pattern.ReplaceAllStringFunc(line, func(match string) string {
		return h.colour + match + Reset
	})
}

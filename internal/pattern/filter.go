package pattern

import "errors"

// LineFilter wraps a Matcher and exposes a Keep method compatible with
// the rest of the logslice pipeline (same contract as filter.TimeFilter).
type LineFilter struct {
	m *Matcher
}

// NewLineFilter constructs a LineFilter from an already-built Matcher.
func NewLineFilter(m *Matcher) (*LineFilter, error) {
	if m == nil {
		return nil, errors.New("pattern: matcher must not be nil")
	}
	return &LineFilter{m: m}, nil
}

// Keep returns true when the line body satisfies the underlying Matcher.
func (f *LineFilter) Keep(line string) bool {
	return f.m.Match(line)
}

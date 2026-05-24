// Package pattern provides substring and regex-based log line matching.
package pattern

import (
	"errors"
	"regexp"
)

// Matcher tests whether a log line body matches a configured pattern.
type Matcher struct {
	re      *regexp.Regexp
	invert  bool
}

// Option configures a Matcher.
type Option func(*Matcher)

// WithInvert causes the matcher to keep lines that do NOT match.
func WithInvert() Option {
	return func(m *Matcher) { m.invert = true }
}

// New compiles pattern as a regular expression and returns a Matcher.
// An empty pattern is not allowed.
func New(pattern string, opts ...Option) (*Matcher, error) {
	if pattern == "" {
		return nil, errors.New("pattern: expression must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	m := &Matcher{re: re}
	for _, o := range opts {
		o(m)
	}
	return m, nil
}

// Match reports whether line satisfies the matcher's rule.
func (m *Matcher) Match(line string) bool {
	matched := m.re.MatchString(line)
	if m.invert {
		return !matched
	}
	return matched
}

// Pattern returns the original regular-expression string.
func (m *Matcher) Pattern() string { return m.re.String() }

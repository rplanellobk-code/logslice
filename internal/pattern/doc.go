// Package pattern implements regex-based log-line filtering for the
// logslice pipeline.
//
// # Overview
//
// A Matcher compiles a regular expression and tests individual log lines
// against it. An optional WithInvert option reverses the match semantics so
// that only lines which do NOT match the expression are retained — useful for
// excluding noisy DEBUG output without modifying other pipeline stages.
//
// A LineFilter wraps a Matcher and exposes a Keep(line string) bool method
// that is compatible with the broader logslice filter contract, making it
// straightforward to compose with TimeFilter and other pipeline components.
//
// # Usage
//
//	m, err := pattern.New(`ERROR|FATAL`)
//	f, err := pattern.NewLineFilter(m)
//	if f.Keep(line) { /* forward line */ }
package pattern

// Package truncate provides line-length truncation for log output.
// Lines exceeding the configured maximum byte length are trimmed and
// optionally suffixed with an ellipsis marker so downstream consumers
// can detect that content was removed.
package truncate

import (
	"errors"
	"fmt"
)

// DefaultSuffix is appended to lines that are truncated.
const DefaultSuffix = "..."

// Truncator trims log lines that exceed a maximum byte length.
type Truncator struct {
	maxLen int
	suffix []byte
}

// New returns a Truncator that caps lines at maxLen bytes.
// If addSuffix is true, DefaultSuffix is appended to truncated lines so
// the total length does not exceed maxLen.
// maxLen must be greater than len(DefaultSuffix) when addSuffix is true.
func New(maxLen int, addSuffix bool) (*Truncator, error) {
	if maxLen <= 0 {
		return nil, errors.New("truncate: maxLen must be positive")
	}
	var suffix []byte
	if addSuffix {
		if maxLen <= len(DefaultSuffix) {
			return nil, fmt.Errorf(
				"truncate: maxLen (%d) must be greater than suffix length (%d)",
				maxLen, len(DefaultSuffix),
			)
		}
		suffix = []byte(DefaultSuffix)
	}
	return &Truncator{maxLen: maxLen, suffix: suffix}, nil
}

// Apply returns a copy of line truncated to at most t.maxLen bytes.
// If the line is within the limit it is returned unchanged.
// When a suffix is configured the last len(suffix) bytes of the allowed
// window are replaced with the suffix so the caller can identify trimmed
// lines without exceeding the budget.
func (t *Truncator) Apply(line []byte) []byte {
	if len(line) <= t.maxLen {
		return line
	}
	out := make([]byte, t.maxLen)
	copy(out, line[:t.maxLen])
	if len(t.suffix) > 0 {
		copy(out[t.maxLen-len(t.suffix):], t.suffix)
	}
	return out
}

// MaxLen returns the configured maximum line length in bytes.
func (t *Truncator) MaxLen() int { return t.maxLen }

// Package sanitize provides line-level sanitization for log output,
// removing or replacing control characters and non-printable bytes
// that could corrupt terminal output or downstream processing.
package sanitize

import (
	"strings"
	"unicode"
)

// Sanitizer scrubs log lines before they are written to output.
type Sanitizer struct {
	replacement rune
	stripANSI   bool
}

// Option configures a Sanitizer.
type Option func(*Sanitizer)

// WithReplacement sets the rune used to replace non-printable characters.
// Defaults to the Unicode replacement character (U+FFFD).
func WithReplacement(r rune) Option {
	return func(s *Sanitizer) { s.replacement = r }
}

// WithStripANSI instructs the sanitizer to remove ANSI escape sequences.
func WithStripANSI(strip bool) Option {
	return func(s *Sanitizer) { s.stripANSI = strip }
}

// New creates a Sanitizer with the supplied options.
func New(opts ...Option) *Sanitizer {
	s := &Sanitizer{replacement: unicode.ReplacementChar}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Apply returns a sanitized copy of line.
// Tab (\t) is preserved; all other non-printable runes are replaced.
func (s *Sanitizer) Apply(line string) string {
	if s.stripANSI {
		line = stripANSIEscapes(line)
	}
	var b strings.Builder
	b.Grow(len(line))
	for _, r := range line {
		if r == '\t' || unicode.IsPrint(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune(s.replacement)
		}
	}
	return b.String()
}

// stripANSIEscapes removes ESC [ … m style ANSI colour/control sequences.
func stripANSIEscapes(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	inEsc := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case inEsc:
			if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
				inEsc = false
			}
		case c == 0x1B && i+1 < len(s) && s[i+1] == '[':
			inEsc = true
			i++ // skip '['
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

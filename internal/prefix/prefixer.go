// Package prefix provides a line transformer that prepends a static or
// dynamic prefix string to every log line's raw text.
package prefix

import (
	"errors"
	"strings"

	"github.com/logslice/logslice/internal/reader"
)

// Prefixer prepends a fixed string to the raw text of each log line.
type Prefixer struct {
	prefix string
}

// Option is a functional option for Prefixer.
type Option func(*Prefixer)

// WithSeparator appends sep between the prefix and the original line text.
// Defaults to a single space when no separator option is supplied.
func WithSeparator(sep string) Option {
	return func(p *Prefixer) {
		// Re-attach separator to the stored prefix so Apply stays simple.
		base := strings.TrimRight(p.prefix, " \t")
		p.prefix = base + sep
	}
}

// New creates a Prefixer that prepends prefix (followed by a space) to every
// line. An empty prefix string is rejected because it would be a no-op.
func New(prefix string, opts ...Option) (*Prefixer, error) {
	if strings.TrimSpace(prefix) == "" {
		return nil, errors.New("prefix: prefix string must not be empty")
	}
	p := &Prefixer{prefix: prefix + " "}
	for _, o := range opts {
		o(p)
	}
	return p, nil
}

// Apply returns a copy of line with the prefix prepended to its Raw field.
func (p *Prefixer) Apply(line reader.LogLine) reader.LogLine {
	line.Raw = p.prefix + line.Raw
	return line
}

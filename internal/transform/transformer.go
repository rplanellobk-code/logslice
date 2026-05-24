// Package transform provides line-level text transformation for log output.
// It applies a chain of user-defined rewrite functions to each log line's
// raw text before it is written to the destination.
package transform

import "github.com/yourorg/logslice/internal/reader"

// Func is a function that transforms a single log line's raw text.
// It receives the original text and returns the (possibly modified) text.
type Func func(text string) string

// Transformer applies an ordered chain of Func values to every LogLine
// it processes.
type Transformer struct {
	fns []Func
}

// New creates a Transformer that will apply fns in order.
// If fns is empty, Apply returns lines unchanged.
func New(fns ...Func) (*Transformer, error) {
	chain := make([]Func, 0, len(fns))
	for _, f := range fns {
		if f == nil {
			return nil, ErrNilFunc
		}
		chain = append(chain, f)
	}
	return &Transformer{fns: chain}, nil
}

// Apply runs the transformation chain against line and returns a new
// reader.LogLine with the updated text. The timestamp is preserved.
func (t *Transformer) Apply(line reader.LogLine) reader.LogLine {
	text := line.Text
	for _, f := range t.fns {
		text = f(text)
	}
	return reader.LogLine{Timestamp: line.Timestamp, Text: text}
}

// ApplyAll transforms every line in src, appending results to dst.
// dst may be nil; a new slice is allocated in that case.
func (t *Transformer) ApplyAll(src []reader.LogLine) []reader.LogLine {
	out := make([]reader.LogLine, len(src))
	for i, l := range src {
		out[i] = t.Apply(l)
	}
	return out
}

package window

import "errors"

// ContextFilter wraps a Buffer and emits a deduplicated set of context
// lines whenever Flush is called. It is useful for capturing the N
// lines that preceded a log event of interest.
type ContextFilter struct {
	buf  *Buffer
	seen map[string]struct{}
}

// NewContextFilter creates a ContextFilter backed by a Buffer of capacity n.
func NewContextFilter(n int) (*ContextFilter, error) {
	if n < 1 {
		return nil, errors.New("window: ContextFilter capacity must be >= 1")
	}
	b, err := New(n)
	if err != nil {
		return nil, err
	}
	return &ContextFilter{
		buf:  b,
		seen: make(map[string]struct{}),
	}, nil
}

// Feed adds a line to the internal rolling buffer.
func (cf *ContextFilter) Feed(line string) {
	cf.buf.Push(line)
}

// Flush returns all currently buffered lines that have not been returned
// by a previous Flush call, preserving chronological order.
// The deduplication set is reset so subsequent Flushes can re-emit new lines.
func (cf *ContextFilter) Flush() []string {
	var out []string
	for _, l := range cf.buf.Lines() {
		if _, dup := cf.seen[l]; !dup {
			out = append(out, l)
			cf.seen[l] = struct{}{}
		}
	}
	// Reset seen so future windows can re-emit lines.
	cf.seen = make(map[string]struct{})
	return out
}

// Reset clears both the buffer and the deduplication set.
func (cf *ContextFilter) Reset() {
	cf.buf.Reset()
	cf.seen = make(map[string]struct{})
}

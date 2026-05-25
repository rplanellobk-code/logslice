// Package window implements a fixed-capacity circular line buffer.
//
// It is intended for use in log processing pipelines where a rolling
// context window of the most-recent N lines is needed — for example,
// to emit surrounding lines when a pattern match is detected.
//
// # Usage
//
//	buf, err := window.New(5)
//	if err != nil { ... }
//
//	for _, line := range incoming {
//	    buf.Push(line)
//	}
//	context := buf.Lines() // up to 5 most-recent lines, oldest first
//
// The buffer allocates once at construction and performs no further
// heap allocations during Push or Lines calls.
package window

// Package tail provides a Tailer that efficiently extracts the last N lines
// from an io.Reader without buffering the entire input into memory.
//
// It uses a fixed-size circular buffer of length N, so memory usage is
// proportional to N rather than to the total input size. This makes it
// suitable for tailing large log archives where only the most recent entries
// are needed.
//
// Basic usage:
//
//	tl, err := tail.New(100)
//	if err != nil { ... }
//	lines, err := tl.Read(f)
//	if err != nil { ... }
//	for _, l := range lines {
//		fmt.Println(l)
//	}
package tail

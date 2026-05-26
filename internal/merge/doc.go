// Package merge implements a k-way merge of pre-sorted log line streams.
//
// When logslice processes a large archive in parallel across multiple goroutines
// or files, each pipeline emits lines sorted by timestamp within its own range.
// The Merger recombines those streams into a single globally-sorted output using
// a min-heap so that the per-line cost is O(log k) where k is the number of
// sources.
//
// Typical usage:
//
//	ch1 := runPipeline(fileA)
//	ch2 := runPipeline(fileB)
//	m, err := merge.New([]<-chan reader.LogLine{ch1, ch2})
//	if err != nil { ... }
//	for line := range m.Merge() {
//	    writer.WriteLine(line)
//	}
package merge

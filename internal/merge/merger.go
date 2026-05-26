// Package merge provides ordered merging of pre-sorted log line streams.
// Lines from multiple sources are merged in ascending timestamp order,
// making it suitable for combining output from parallel pipeline runs.
package merge

import (
	"container/heap"
	"errors"

	"github.com/logslice/logslice/internal/reader"
)

// ErrNoSources is returned when New is called with an empty source slice.
var ErrNoSources = errors.New("merge: at least one source is required")

// entry holds a buffered line and the index of the source it came from.
type entry struct {
	line   reader.LogLine
	source int
}

// minHeap implements heap.Interface for entry values ordered by timestamp.
type minHeap []entry

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].line.Timestamp.Before(h[j].line.Timestamp) }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(entry)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Merger merges multiple sorted LogLine channels into a single sorted stream.
type Merger struct {
	sources []<-chan reader.LogLine
}

// New creates a Merger from the given source channels.
// Each source must produce lines in ascending timestamp order.
func New(sources []<-chan reader.LogLine) (*Merger, error) {
	if len(sources) == 0 {
		return nil, ErrNoSources
	}
	return &Merger{sources: sources}, nil
}

// Merge reads from all sources and emits lines in ascending timestamp order.
// The returned channel is closed once all sources are exhausted.
func (m *Merger) Merge() <-chan reader.LogLine {
	out := make(chan reader.LogLine, len(m.sources)*4)
	go func() {
		defer close(out)
		h := &minHeap{}
		heap.Init(h)
		for i, src := range m.sources {
			if line, ok := <-src; ok {
				heap.Push(h, entry{line: line, source: i})
			}
		}
		for h.Len() > 0 {
			e := heap.Pop(h).(entry)
			out <- e.line
			if line, ok := <-m.sources[e.source]; ok {
				heap.Push(h, entry{line: line, source: e.source})
			}
		}
	}()
	return out
}

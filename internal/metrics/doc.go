// Package metrics provides lightweight, goroutine-safe counters for tracking
// the throughput of a logslice pipeline run.
//
// # Overview
//
// A [Counters] value is created with [New] at the start of a pipeline run.
// Each stage of the pipeline increments the relevant counter atomically:
//
//	- LinesRead    — every line consumed from a source file.
//	- LinesMatched — lines whose timestamp falls within the requested range.
//	- LinesWritten — lines successfully flushed to the output writer.
//	- LinesSkipped — lines discarded (parse error, out-of-range, etc.).
//	- BytesWritten — raw bytes handed to the output writer.
//
// At the end of a run, [Counters.Summary] prints a human-readable report and
// [Counters.Rate] exposes lines-per-second throughput for progress displays.
//
// # Thread Safety
//
// All counter fields use [sync/atomic.Int64] internally, so they may be
// incremented concurrently from multiple goroutines without additional locking.
package metrics

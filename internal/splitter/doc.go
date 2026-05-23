// Package splitter provides the core orchestration layer for logslice.
//
// A Splitter ties together a LineReader, a TimeFilter, and a LineWriter to
// extract log lines that fall within a requested time window.  Progress
// tracking is optional: when a non-nil *progress.Tracker is supplied in the
// Config, every line read (regardless of whether it passes the filter) is
// counted so that callers can report throughput and ETA to the user.
//
// Typical usage:
//
//	 sp, err := splitter.New(splitter.Config{
//	     Source:  sourceReader,
//	     Dest:    destWriter,
//	     Filter:  timeFilter,
//	     Tracker: tracker, // optional
//	 })
//	 if err != nil { ... }
//	 written, err := sp.Run(lineReader, lineWriter)
package splitter

// Package dedupe implements a sliding-window deduplication filter for
// log line streams.
//
// # Overview
//
// When processing large log archives it is common to encounter runs of
// identical lines (e.g. repeated error messages or heartbeat entries).
// The Filter type tracks the last N unique lines seen and reports
// whether a new line is a duplicate of one already in that window.
//
// # Usage
//
//	f, err := dedupe.New(1000)
//	if err != nil { ... }
//	for _, line := range lines {
//		if !f.IsDuplicate(line) {
//			// emit line
//		}
//	}
//
// The window capacity is configurable so callers can balance memory
// usage against the breadth of deduplication.
package dedupe

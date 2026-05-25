// Package burst provides a sliding-window burst detector for log line streams.
//
// A Detector maintains a ring of recent arrival timestamps and reports whether
// the volume of lines within a configurable time window has exceeded a
// configured threshold. This is useful for alerting or back-pressure when a
// log source emits an unexpectedly high number of lines in a short period.
//
// Basic usage:
//
//	det, err := burst.New(5*time.Second, 1000)
//	if err != nil { /* handle */ }
//
//	for _, line := range lines {
//		if det.Record(line.Timestamp) {
//			log.Println("burst detected")
//		}
//	}
package burst

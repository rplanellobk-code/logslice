// Package progress provides a lightweight Tracker that counts bytes and lines
// processed during a logslice run and periodically reports statistics to any
// io.Writer (typically os.Stderr).
//
// Usage:
//
//	tr := progress.NewTracker(os.Stderr, 5*time.Second)
//	defer tr.Stop()
//
//	// inside your read loop:
//	tr.AddBytes(int64(len(line)))
//	tr.AddLineIn()
//
//	// after the filter accepts a line:
//	tr.AddLineOut()
package progress

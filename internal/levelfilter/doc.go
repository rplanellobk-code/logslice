// Package levelfilter implements severity-level-based filtering for log lines.
//
// A Filter is constructed with a minimum level string such as "warn" or
// "error". Calling Keep on a raw log line returns true only when the line
// contains a recognised severity keyword that is at or above the configured
// minimum.
//
// Supported level tokens (case-insensitive):
//
//	debug
//	info
//	warn / warning
//	error / err
//	fatal / crit
//
// Example:
//
//	f, err := levelfilter.New("warn")
//	if err != nil { ... }
//	if f.Keep(line) {
//	    // write line to output
//	}
package levelfilter

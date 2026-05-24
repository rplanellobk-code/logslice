// Package rotate implements output file rotation for logslice.
//
// A Rotator opens a new numbered output file whenever the current file
// exceeds a configurable threshold. Two rotation strategies are supported:
//
//   - ByLines: rotate after a fixed number of log lines have been written.
//   - BySize: rotate after the cumulative byte count of the current file
//     would exceed the threshold.
//
// Output files are named <prefix>_NNNN.log (zero-padded to four digits) and
// are created inside the directory supplied to New.
//
// Example usage:
//
//	r, err := rotate.New("/var/log/out", "app", rotate.ByLines, 10_000)
//	if err != nil { ... }
//	defer r.Close()
//	for _, line := range lines {
//		_ = r.WriteLine(line)
//	}
package rotate

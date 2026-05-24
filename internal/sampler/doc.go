// Package sampler implements deterministic 1-in-N log line sampling.
//
// When processing very large log archives it is often useful to reduce
// the output to a representative subset of lines rather than every
// matching entry. Sampler provides a thread-safe counter-based approach
// that keeps the first line of every N-line window, making the output
// evenly distributed across the original time range.
//
// Basic usage:
//
//	s, err := sampler.New(10) // keep every 10th line
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range lines {
//		if s.Keep(line) {
//			output(line)
//		}
//	}
package sampler

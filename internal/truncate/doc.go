// Package truncate implements a byte-length limiter for individual log lines.
//
// When processing large log archives it is common to encounter lines that are
// pathologically long (e.g. base64-encoded blobs or minified JSON). Such lines
// can cause downstream consumers to allocate excessive memory or exceed column
// limits in storage systems.
//
// Usage:
//
//	tr, err := truncate.New(512, true)
//	if err != nil {
//		log.Fatal(err)
//	}
//	processed := tr.Apply(rawLine)
//
// When addSuffix is true the last three bytes of the truncated output are
// replaced with "..." so readers can distinguish trimmed lines from lines
// that happen to end at the limit.
package truncate

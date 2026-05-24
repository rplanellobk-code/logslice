// Package mask provides regex-based redaction of sensitive data within log
// lines. A Masker is constructed with one or more compiled regular expressions
// and an optional replacement string. Calling Apply on a raw log line returns
// a copy of the line with every pattern match replaced by the mask token,
// making it safe to write personal or secret data to untrusted destinations.
//
// Usage:
//
//	masker, err := mask.New(
//		[]string{`password=\S+`, `token=[A-Za-z0-9]+`},
//		mask.WithMask("[REDACTED]"),
//	)
//	if err != nil { /* handle */ }
//	clean := masker.Apply(rawLine)
package mask

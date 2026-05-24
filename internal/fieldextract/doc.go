// Package fieldextract provides lightweight key/value field extraction
// for structured log lines.
//
// # Overview
//
// Log lines often carry metadata as space-separated key=value (or key:value)
// pairs, for example:
//
//	ts=2024-01-01T12:00:00Z level=info svc=api code=200 msg="request ok"
//
// Extractor locates a single named field within such a line and returns its
// value. MultiExtractor wraps several Extractor instances so that multiple
// fields can be harvested from a line in a single call.
//
// # Usage
//
//	ex, _ := fieldextract.New("level", "=")
//	val, ok := ex.Extract(line)
//
//	mex, _ := fieldextract.NewMulti(map[string]string{"level": "=", "code": "="})
//	fields := mex.ExtractAll(line)
package fieldextract

// Package index builds and queries a lightweight byte-offset index over
// a structured log file.
//
// # Overview
//
// Rather than scanning an entire file to locate a time range, Build walks
// the file once and records the byte offset of every line whose timestamp
// can be parsed. The resulting FileIndex can then answer two queries in
// O(n) time:
//
//   - FindStart — the offset of the first line whose timestamp is >= a
//     given lower bound.
//   - FindEnd — the offset of the last line whose timestamp is <= a given
//     upper bound.
//
// These offsets can be handed directly to io.ReadSeeker.Seek, allowing the
// caller to skip straight to the relevant portion of the file.
//
// # Usage
//
//	idx, err := index.Build(file, parser)
//	start := idx.FindStart(from)
//	file.Seek(start, io.SeekStart)
package index

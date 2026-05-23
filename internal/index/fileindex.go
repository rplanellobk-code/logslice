// Package index provides byte-offset indexing for log files,
// enabling fast seeking to time-range boundaries without full scans.
package index

import (
	"bufio"
	"io"
	"time"

	"github.com/user/logslice/internal/timestamp"
)

// Entry records the byte offset and parsed timestamp of a log line.
type Entry struct {
	Offset    int64
	Timestamp time.Time
}

// FileIndex holds an ordered slice of index entries for a log file.
type FileIndex struct {
	Entries []Entry
}

// Build scans r, parsing timestamps with p, and returns a FileIndex
// containing one Entry per successfully parsed line.
func Build(r io.ReadSeeker, p *timestamp.Parser) (*FileIndex, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var (
		idx    FileIndex
		offset int64
		scanner = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		line := scanner.Text()
		if ts, err := p.Parse(line); err == nil {
			idx.Entries = append(idx.Entries, Entry{
				Offset:    offset,
				Timestamp: ts,
			})
		}
		// +1 for the newline byte consumed by the scanner
		offset += int64(len(line)) + 1
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &idx, nil
}

// FindStart returns the byte offset of the first entry whose timestamp
// is >= start. Returns 0 if no bound is needed (zero value of time.Time).
func (idx *FileIndex) FindStart(start time.Time) int64 {
	if start.IsZero() || len(idx.Entries) == 0 {
		return 0
	}
	for _, e := range idx.Entries {
		if !e.Timestamp.Before(start) {
			return e.Offset
		}
	}
	return -1 // all entries are before start
}

// FindEnd returns the byte offset just after the last entry whose
// timestamp is <= end. Returns -1 if no bound is needed.
func (idx *FileIndex) FindEnd(end time.Time) int64 {
	if end.IsZero() || len(idx.Entries) == 0 {
		return -1
	}
	var last int64 = -1
	for _, e := range idx.Entries {
		if !e.Timestamp.After(end) {
			last = e.Offset
		}
	}
	return last
}

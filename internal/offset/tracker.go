// Package offset tracks byte offsets within a log file, enabling
// efficient seek operations and resume-from-position functionality.
package offset

import (
	"errors"
	"io"
	"sync/atomic"
)

// ErrNegativeOffset is returned when a negative offset is provided.
var ErrNegativeOffset = errors.New("offset: negative offset")

// Tracker wraps an io.ReadSeeker and records the current byte offset
// as lines are consumed.
type Tracker struct {
	rs      io.ReadSeeker
	offset  atomic.Int64
	buf     []byte
}

// New creates a Tracker wrapping rs. The initial offset is set to
// the current position of rs (typically 0).
func New(rs io.ReadSeeker) (*Tracker, error) {
	if rs == nil {
		return nil, errors.New("offset: nil ReadSeeker")
	}
	pos, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	t := &Tracker{
		rs:  rs,
		buf: make([]byte, 4096),
	}
	t.offset.Store(pos)
	return t, nil
}

// ReadLine reads the next newline-terminated line from the underlying
// reader, advancing the tracked offset accordingly.
// Returns io.EOF when no more data is available.
func (t *Tracker) ReadLine() (string, error) {
	var line []byte
	for {
		n, err := t.rs.Read(t.buf[:1])
		if n == 1 {
			t.offset.Add(1)
			if t.buf[0] == '\n' {
				return string(line), nil
			}
			line = append(line, t.buf[0])
		}
		if err != nil {
			if err == io.EOF && len(line) > 0 {
				return string(line), nil
			}
			return "", err
		}
	}
}

// Offset returns the current byte offset within the underlying reader.
func (t *Tracker) Offset() int64 {
	return t.offset.Load()
}

// SeekTo moves the underlying reader to the given byte offset and
// updates the tracked position.
func (t *Tracker) SeekTo(off int64) error {
	if off < 0 {
		return ErrNegativeOffset
	}
	pos, err := t.rs.Seek(off, io.SeekStart)
	if err != nil {
		return err
	}
	t.offset.Store(pos)
	return nil
}

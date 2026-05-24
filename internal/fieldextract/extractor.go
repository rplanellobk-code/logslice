// Package fieldextract provides utilities for extracting named fields
// from structured log lines (e.g. key=value or JSON-style pairs).
package fieldextract

import (
	"fmt"
	"strings"
)

// Extractor extracts a named field value from a raw log line.
type Extractor struct {
	field     string
	delimiter string
}

// New returns an Extractor that looks for the given field name using
// the provided key/value delimiter (e.g. "=" or ":").
// Returns an error if field or delimiter is empty.
func New(field, delimiter string) (*Extractor, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldextract: field name must not be empty")
	}
	if delimiter == "" {
		return nil, fmt.Errorf("fieldextract: delimiter must not be empty")
	}
	return &Extractor{field: field, delimiter: delimiter}, nil
}

// Extract scans line for a token of the form "<field><delimiter><value>" and
// returns the value. The value ends at the next whitespace character or end
// of string. Returns ("", false) when the field is not found.
func (e *Extractor) Extract(line string) (string, bool) {
	prefix := e.field + e.delimiter
	idx := strings.Index(line, prefix)
	if idx == -1 {
		return "", false
	}
	start := idx + len(prefix)
	if start >= len(line) {
		return "", false
	}
	rest := line[start:]
	end := strings.IndexAny(rest, " \t\r\n")
	if end == -1 {
		return rest, true
	}
	value := rest[:end]
	if value == "" {
		return "", false
	}
	return value, true
}

// Field returns the configured field name.
func (e *Extractor) Field() string { return e.field }

// Delimiter returns the configured delimiter.
func (e *Extractor) Delimiter() string { return e.delimiter }

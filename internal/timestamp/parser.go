package timestamp

import (
	"fmt"
	"time"
)

// Common log timestamp formats to try when parsing.
var knownFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05.000000",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
	"Jan  2 15:04:05",
}

// Parser attempts to extract and parse a timestamp from a log line.
type Parser struct {
	format string
	loc    *time.Location
}

// NewParser creates a Parser. If format is empty, auto-detection is used.
func NewParser(format string, loc *time.Location) *Parser {
	if loc == nil {
		loc = time.UTC
	}
	return &Parser{format: format, loc: loc}
}

// Parse attempts to parse a timestamp from the beginning of a log line.
// It returns the parsed time and the number of bytes consumed, or an error.
func (p *Parser) Parse(line string) (time.Time, error) {
	if p.format != "" {
		return p.parseWith(line, p.format)
	}
	return p.autoDetect(line)
}

func (p *Parser) parseWith(line, format string) (time.Time, error) {
	// Try to parse from the start, allowing for a prefix scan up to 40 chars.
	maxPrefix := len(format) + 10
	if maxPrefix > len(line) {
		maxPrefix = len(line)
	}
	t, err := time.ParseInLocation(format, line[:maxPrefix], p.loc)
	if err == nil {
		return t, nil
	}
	// Try exact slice matching format length
	if len(line) >= len(format) {
		t, err = time.ParseInLocation(format, line[:len(format)], p.loc)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("timestamp: cannot parse %q with format %q", line, format)
}

func (p *Parser) autoDetect(line string) (time.Time, error) {
	for _, fmt := range knownFormats {
		if t, err := p.parseWith(line, fmt); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("timestamp: unable to detect timestamp in line: %.60s", line)
}

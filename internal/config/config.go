// Package config holds the validated runtime configuration for logslice.
package config

import (
	"errors"
	"time"
)

// Config holds all user-supplied options for a logslice run.
type Config struct {
	// Pattern is a glob or directory path that identifies input log files.
	Pattern string

	// OutputDir is the directory where extracted segments are written.
	OutputDir string

	// From is the inclusive start of the time range filter (zero = open).
	From time.Time

	// To is the exclusive end of the time range filter (zero = open).
	To time.Time

	// TimestampFormat is an optional explicit Go time layout string.
	// When empty the parser auto-detects the format.
	TimestampFormat string

	// SampleN keeps every Nth matching line (0 or 1 = keep all).
	SampleN uint64

	// MaxLinesPerFile rotates the output file after this many lines (0 = no rotation).
	MaxLinesPerFile uint64

	// MaxBytesPerFile rotates the output file after this many bytes (0 = no rotation).
	MaxBytesPerFile uint64

	// NormalizeTimestamps rewrites timestamps to a canonical format on output.
	NormalizeTimestamps bool

	// CheckpointPath is an optional file path for resumable processing state.
	CheckpointPath string

	// DedupeWindow is the number of recent lines checked for duplicates (0 = off).
	DedupeWindow int

	// RateLimit is the maximum lines-per-second written (0 = unlimited).
	RateLimit float64
}

// Validate returns an error if the configuration is not usable.
func Validate(c Config) error {
	if c.Pattern == "" {
		return errors.New("config: pattern must not be empty")
	}
	if c.OutputDir == "" {
		return errors.New("config: output_dir must not be empty")
	}
	if !c.From.IsZero() && !c.To.IsZero() && !c.To.After(c.From) {
		return errors.New("config: to must be after from")
	}
	if c.DedupeWindow < 0 {
		return errors.New("config: dedupe_window must be >= 0")
	}
	if c.RateLimit < 0 {
		return errors.New("config: rate_limit must be >= 0")
	}
	return nil
}

// Package config provides configuration loading and validation for logslice.
package config

import (
	"errors"
	"fmt"
	"time"
)

// Config holds the runtime configuration for a logslice operation.
type Config struct {
	// Pattern is a file path or glob pattern pointing to log files.
	Pattern string

	// From is the inclusive start of the time range filter. Zero means open.
	From time.Time

	// To is the inclusive end of the time range filter. Zero means open.
	To time.Time

	// TimestampFormat is an optional explicit strftime/Go layout for parsing
	// timestamps. When empty, auto-detection is used.
	TimestampFormat string

	// OutputDir is the directory where output files are written.
	OutputDir string

	// NormalizeTimestamps rewrites timestamps in output lines to RFC3339.
	NormalizeTimestamps bool

	// Quiet suppresses progress output.
	Quiet bool
}

// Validate checks that the Config fields are self-consistent.
// It returns a descriptive error for the first violation found.
func (c *Config) Validate() error {
	if c.Pattern == "" {
		return errors.New("config: pattern must not be empty")
	}
	if c.OutputDir == "" {
		return errors.New("config: output directory must not be empty")
	}
	if !c.From.IsZero() && !c.To.IsZero() && c.To.Before(c.From) {
		return fmt.Errorf(
			"config: 'to' time %s is before 'from' time %s",
			c.To.Format(time.RFC3339),
			c.From.Format(time.RFC3339),
		)
	}
	return nil
}

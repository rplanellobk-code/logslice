// logslice is a fast log file splitter and time-range extractor
// for large structured log archives.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourusername/logslice/internal/filter"
	"github.com/yourusername/logslice/internal/reader"
	"github.com/yourusername/logslice/internal/timestamp"
	"github.com/yourusername/logslice/internal/writer"
)

const usageText = `logslice — extract a time range from a structured log file.

Usage:
  logslice [flags] [input-file]

If no input file is given, logslice reads from stdin.

Flags:
`

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprint(stderr, usageText)
		fs.PrintDefaults()
	}

	var (
		sinceStr  = fs.String("since", "", "include lines at or after this timestamp (RFC3339 or common log formats)")
		untilStr  = fs.String("until", "", "include lines before or at this timestamp")
		formatStr = fs.String("format", "", "explicit Go time layout for parsing log timestamps (optional)")
		timezone  = fs.String("tz", "UTC", "timezone to use when parsing timestamps without zone info")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Resolve timezone.
	loc, err := time.LoadLocation(*timezone)
	if err != nil {
		return fmt.Errorf("invalid timezone %q: %w", *timezone, err)
	}

	// Build timestamp parser.
	parser, err := timestamp.NewParser(*formatStr, loc)
	if err != nil {
		return fmt.Errorf("timestamp parser: %w", err)
	}

	// Parse optional since/until bounds.
	var since, until time.Time
	if *sinceStr != "" {
		since, err = parser.Parse(*sinceStr)
		if err != nil {
			return fmt.Errorf("--since: %w", err)
		}
	}
	if *untilStr != "" {
		until, err = parser.Parse(*untilStr)
		if err != nil {
			return fmt.Errorf("--until: %w", err)
		}
	}

	// Build the time filter.
	tf, err := filter.NewTimeFilter(since, until)
	if err != nil {
		return fmt.Errorf("time filter: %w", err)
	}

	// Open input source.
	var input io.Reader = stdin
	if fs.NArg() > 0 {
		f, err := os.Open(fs.Arg(0))
		if err != nil {
			return fmt.Errorf("open input: %w", err)
		}
		defer f.Close()
		input = f
	}

	// Wire up the pipeline: reader → filter → writer.
	lr := reader.NewLineReader(input, parser)
	lw := writer.NewLineWriter(stdout)

	for {
		line, err := lr.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Log parse errors to stderr and continue; don't abort the run.
			fmt.Fprintf(stderr, "warning: %v\n", err)
			continue
		}

		if !tf.Include(line) {
			continue
		}

		if err := lw.WriteLine(line); err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}

	return nil
}

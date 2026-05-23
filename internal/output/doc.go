// Package output provides formatting primitives for rendering extracted log
// lines to an io.Writer.
//
// # Formats
//
// Two output formats are supported:
//
//   - FormatRaw – writes each log line exactly as it appeared in the source
//     file, followed by a newline.
//
//   - FormatNormalized – prepends a re-formatted timestamp (defaulting to
//     RFC3339Nano) to the original raw text, separated by a tab. This is
//     useful when combining logs from multiple sources that use different
//     timestamp dialects.
//
// # Usage
//
//	f, err := output.NewFormatter(os.Stdout, output.FormatRaw, "")
//	if err != nil { ... }
//	for _, line := range lines {
//	    f.Write(line)
//	}
package output

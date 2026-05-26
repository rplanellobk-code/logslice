package highlight

import "fmt"

// LineHighlighter adapts Highlighter to a simple string-transform function
// compatible with the transform.Func signature used in the pipeline.
//
// Usage:
//
//	step := highlight.PipelineStep(`ERROR`, highlight.WithColour(highlight.Red))
//	transformer, err := transform.New(step)
func PipelineStep(pattern string, opts ...Option) (func(string) (string, error), error) {
	h, err := New(pattern, opts...)
	if err != nil {
		return nil, fmt.Errorf("highlight.PipelineStep: %w", err)
	}
	return func(line string) (string, error) {
		return h.Apply(line), nil
	}, nil
}

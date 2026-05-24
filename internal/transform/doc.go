// Package transform provides composable, line-level text transformations
// for log processing pipelines.
//
// A Transformer holds an ordered chain of Func values. When Apply or
// ApplyAll is called each function in the chain receives the output of
// the previous one, enabling flexible rewrite pipelines such as:
//
//   - field redaction (e.g. masking IP addresses or tokens)
//   - prefix / suffix injection
//   - whitespace normalisation
//   - custom codec translation
//
// Usage:
//
//	tr, err := transform.New(
//	    strings.ToUpper,
//	    func(s string) string { return "[LOG] " + s },
//	)
//	if err != nil { … }
//	outLine := tr.Apply(inLine)
package transform

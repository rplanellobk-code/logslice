package transform

import "errors"

// ErrNilFunc is returned by New when a nil Func is supplied in the chain.
var ErrNilFunc = errors.New("transform: nil Func in chain")

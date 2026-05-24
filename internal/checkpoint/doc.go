// Package checkpoint provides lightweight, file-backed checkpoint persistence
// for logslice processing runs.
//
// A checkpoint captures the last successfully processed file path, byte offset,
// and log-line timestamp so that a subsequent run can skip already-processed
// data and resume from exactly where the previous run left off.
//
// Usage:
//
//	store, err := checkpoint.NewStore("/var/run/logslice/state.json")
//	if err != nil { ... }
//
//	st, err := store.Load()   // returns zero State if file absent
//	if err != nil { ... }
//
//	// ... process logs, update st ...
//
//	if err := store.Save(st); err != nil { ... }
//
// State is serialised as indented JSON and written atomically via a
// temporary-file rename so that a crash mid-write never leaves a corrupt
// checkpoint on disk.
package checkpoint

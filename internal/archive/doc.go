// Package archive provides file discovery for log archives.
//
// It supports two discovery modes:
//
//   - Directory mode: when the supplied path is an existing directory,
//     all regular files within that directory (non-recursive) are returned
//     sorted lexicographically by path.
//
//   - Glob mode: when the supplied path contains glob meta-characters
//     (e.g. "/var/log/app-*.log"), filepath.Glob is used to enumerate
//     matching regular files.
//
// Typical usage:
//
//	s, err := archive.NewScanner("/var/log/myapp")
//	if err != nil { ... }
//	files, err := s.Scan()
//	for _, f := range files {
//		fmt.Println(f.Path, f.Size)
//	}
package archive

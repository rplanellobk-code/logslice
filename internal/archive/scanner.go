// Package archive provides utilities for scanning and enumerating
// log archive files from a directory or glob pattern.
package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// File represents a single log file discovered during scanning.
type File struct {
	Path string
	Size int64
}

// Scanner discovers log files matching a glob pattern or directory.
type Scanner struct {
	pattern string
}

// NewScanner creates a Scanner for the given glob pattern or directory path.
// If pattern is a directory, all files within it are enumerated.
func NewScanner(pattern string) (*Scanner, error) {
	if pattern == "" {
		return nil, fmt.Errorf("archive: pattern must not be empty")
	}
	return &Scanner{pattern: pattern}, nil
}

// Scan returns all matching files sorted by path.
func (s *Scanner) Scan() ([]File, error) {
	info, err := os.Stat(s.pattern)
	if err == nil && info.IsDir() {
		return s.scanDir(s.pattern)
	}
	return s.scanGlob(s.pattern)
}

func (s *Scanner) scanDir(dir string) ([]File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("archive: reading directory %q: %w", dir, err)
	}
	var files []File
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fi, err := e.Info()
		if err != nil {
			return nil, fmt.Errorf("archive: stat %q: %w", e.Name(), err)
		}
		files = append(files, File{
			Path: filepath.Join(dir, e.Name()),
			Size: fi.Size(),
		})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

func (s *Scanner) scanGlob(pattern string) ([]File, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("archive: invalid glob pattern %q: %w", pattern, err)
	}
	var files []File
	for _, m := range matches {
		fi, err := os.Stat(m)
		if err != nil || fi.IsDir() {
			continue
		}
		files = append(files, File{Path: m, Size: fi.Size()})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

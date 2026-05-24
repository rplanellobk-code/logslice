// Package checkpoint provides persistent progress tracking for log processing
// runs, allowing resumption from the last successfully processed position.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// State holds the persisted state of a processing run.
type State struct {
	// LastFile is the path of the last successfully processed file.
	LastFile string `json:"last_file"`
	// LastOffset is the byte offset within LastFile up to which lines were processed.
	LastOffset int64 `json:"last_offset"`
	// LastTimestamp is the timestamp of the last processed log line.
	LastTimestamp time.Time `json:"last_timestamp"`
	// UpdatedAt records when this checkpoint was written.
	UpdatedAt time.Time `json:"updated_at"`
}

// Store manages reading and writing checkpoint state to a file.
type Store struct {
	path string
}

// NewStore creates a Store that persists state at the given file path.
func NewStore(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("checkpoint: path must not be empty")
	}
	return &Store{path: path}, nil
}

// Load reads the checkpoint state from disk. If the file does not exist,
// a zero-value State and nil error are returned.
func (s *Store) Load() (State, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return State{}, nil
		}
		return State{}, err
	}
	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return State{}, err
	}
	return st, nil
}

// Save writes the given state to disk atomically via a temp file rename.
func (s *Store) Save(st State) error {
	st.UpdatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// Delete removes the checkpoint file. If the file does not exist the
// call is a no-op.
func (s *Store) Delete() error {
	err := os.Remove(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

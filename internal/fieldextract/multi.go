package fieldextract

import "fmt"

// MultiExtractor holds several Extractor instances and can extract
// multiple fields from a single log line in one pass.
type MultiExtractor struct {
	extractors []*Extractor
}

// NewMulti builds a MultiExtractor from a map of field→delimiter pairs.
// Returns an error if any field or delimiter is invalid.
func NewMulti(fields map[string]string) (*MultiExtractor, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("fieldextract: at least one field required")
	}
	m := &MultiExtractor{}
	for field, delim := range fields {
		ex, err := New(field, delim)
		if err != nil {
			return nil, err
		}
		m.extractors = append(m.extractors, ex)
	}
	return m, nil
}

// ExtractAll returns a map of field→value for every configured field
// that is present in line. Fields not found are omitted from the map.
func (m *MultiExtractor) ExtractAll(line string) map[string]string {
	out := make(map[string]string, len(m.extractors))
	for _, ex := range m.extractors {
		if val, ok := ex.Extract(line); ok {
			out[ex.Field()] = val
		}
	}
	return out
}

// Fields returns the list of configured field names.
func (m *MultiExtractor) Fields() []string {
	names := make([]string, len(m.extractors))
	for i, ex := range m.extractors {
		names[i] = ex.Field()
	}
	return names
}

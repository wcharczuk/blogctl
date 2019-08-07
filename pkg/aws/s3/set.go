package s3

import "sync"

// Set is a set that can be accessed concurrently.
type Set struct {
	syncroot sync.Mutex
	values   map[string]struct{}
}

// Set adds a value
func (s *Set) Set(value string) {
	s.syncroot.Lock()
	if s.values == nil {
		s.values = make(map[string]struct{})
	}
	s.values[value] = struct{}{}
	s.syncroot.Unlock()
}

// Has returns if a value is in the set
func (s *Set) Has(value string) (has bool) {
	s.syncroot.Lock()
	if s.values == nil {
		s.syncroot.Unlock()
		return
	}
	_, has = s.values[value]
	s.syncroot.Unlock()
	return
}

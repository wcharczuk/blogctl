package web

import "sync"

// Any is an alias to the empty interface.
type Any = interface{}

// Values is an alias to map[string]interface{}
type Values = map[string]interface{}

// State is a provider for a state bag
type State interface {
	Keys() []string
	Get(key string) Any
	Set(key string, value Any)
	Remove(key string)
	Copy() State
}

// SyncState is the collection of state objects on a context.
type SyncState struct {
	sync.Mutex
	Values map[string]Any
}

// Keys returns
func (s *SyncState) Keys() (output []string) {
	s.Lock()
	defer s.Unlock()
	if s.Values == nil {
		return
	}

	output = make([]string, len(s.Values))
	var index int
	for key := range s.Values {
		output[index] = key
		index++
	}
	return
}

// Get gets a value.
func (s *SyncState) Get(key string) Any {
	s.Lock()
	defer s.Unlock()
	if s.Values == nil {
		return nil
	}
	return s.Values[key]
}

// Set sets a value.
func (s *SyncState) Set(key string, value Any) {
	s.Lock()
	defer s.Unlock()
	if s.Values == nil {
		s.Values = make(map[string]Any)
	}
	s.Values[key] = value
	return
}

// Remove removes a key.
func (s *SyncState) Remove(key string) {
	s.Lock()
	defer s.Unlock()
	if s.Values == nil {
		return
	}
	delete(s.Values, key)
	return
}

// Copy creates a new copy of the vars.
func (s *SyncState) Copy() State {
	s.Lock()
	defer s.Unlock()
	return &SyncState{
		Values: s.Values,
	}
}

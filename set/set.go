package set

import "sync"

type Set struct {
	m  map[string]struct{}
	mu sync.Mutex
}

func New() *Set {
	return &Set{
		m: make(map[string]struct{}),
	}
}

func (s *Set) Insert(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = struct{}{}
}

func (s *Set) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
}

func (s *Set) Contains(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.m[key]; ok {
		return true
	}
	return false
}

func (s *Set) Keys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var (
		keys = make([]string, len(s.m))
		i    = 0
	)
	for key := range s.m {
		keys[i] = key
		i++
	}
	return keys
}

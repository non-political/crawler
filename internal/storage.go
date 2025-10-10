package internal

import (
	"sync"
)

type ThreadSafeSet struct {
	mu    sync.RWMutex
	items map[string]struct{}
}

var (
	once         sync.Once
	singletonSet *ThreadSafeSet
)

func NewThreadSafeSet() *ThreadSafeSet {
	return &ThreadSafeSet{
		items: make(map[string]struct{}),
	}
}

// returns the shared set
func Set() *ThreadSafeSet {
	once.Do(func() {
		singletonSet = NewThreadSafeSet()
	})
	return singletonSet
}

func (s *ThreadSafeSet) Add(item string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item] = struct{}{}
}

func (s *ThreadSafeSet) Contains(item string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.items[item]
	return exists
}

func (s *ThreadSafeSet) Delete(item string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, item)
}

func (s *ThreadSafeSet) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

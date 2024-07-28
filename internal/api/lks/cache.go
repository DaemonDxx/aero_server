package lks

import "sync"

type SimpleCookieCache struct {
	c  map[string]map[string]string
	mu sync.RWMutex
}

func NewSimpleCookieCache() *SimpleCookieCache {
	return &SimpleCookieCache{
		c:  make(map[string]map[string]string),
		mu: sync.RWMutex{},
	}
}

func (s *SimpleCookieCache) Get(key string) (map[string]string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.c[key]
	return t, ok
}

func (s *SimpleCookieCache) Put(key string, cookie map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.c[key] = cookie
	return nil
}

func (s *SimpleCookieCache) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.c, key)
}

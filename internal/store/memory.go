package store

import (
	"sync"

	"pingpulse/internal/monitor"
)

type MemoryStore struct {
	mu     sync.RWMutex
	cap    int
	data   map[string][]monitor.Result
}

func NewMemoryStore(capacity int) *MemoryStore {
	return &MemoryStore{
		cap:  capacity,
		data: make(map[string][]monitor.Result),
	}
}

func (s *MemoryStore) Append(r monitor.Result) {
	s.mu.Lock()
	defer s.mu.Unlock()
	buf := append(s.data[r.URL], r)
	if len(buf) > s.cap {
		buf = buf[len(buf)-s.cap:]
	}
	s.data[r.URL] = buf
}

func (s *MemoryStore) Latest(url string, n int) []monitor.Result {
	s.mu.RLock()
	defer s.mu.RUnlock()
	buf := s.data[url]
	if n <= 0 || n > len(buf) {
		n = len(buf)
	}
	start := len(buf) - n
	return append([]monitor.Result(nil), buf[start:]...)
}

func (s *MemoryStore) All() map[string][]monitor.Result {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string][]monitor.Result, len(s.data))
	for k, v := range s.data {
		cp := make([]monitor.Result, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}
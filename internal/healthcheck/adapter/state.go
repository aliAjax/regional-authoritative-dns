package adapter

import (
	"github.com/example/regional-authoritative-dns/internal/healthcheck/domain"
	"sync"
	"time"
)

type State struct {
	mu    sync.RWMutex
	items map[string]domain.Result
}

func New() *State                    { return &State{items: map[string]domain.Result{}} }
func (s *State) Put(r domain.Result) { s.mu.Lock(); defer s.mu.Unlock(); s.items[r.Name] = r }
func (s *State) Get(name string) (domain.Result, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.items[name]
	return r, ok
}
func (s *State) Expired(name string, d time.Duration) bool {
	r, ok := s.Get(name)
	if !ok {
		return true
	}
	// A result whose CheckedAt is in the future (clock skew, or a refresh that
	// stamped ahead of now) is treated as fresh rather than expired, so the
	// negative age must not be reported as expired.
	age := time.Since(r.CheckedAt)
	if age < 0 {
		return false
	}
	return age > d
}

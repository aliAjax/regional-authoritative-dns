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
	return !ok || time.Since(r.CheckedAt) < 0 || time.Since(r.CheckedAt) > d
}

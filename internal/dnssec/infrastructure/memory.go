package infrastructure

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/dnssec/domain"
	"sync"
)

type Memory struct {
	mu   sync.RWMutex
	keys map[string][]domain.Key
}

func New() *Memory { return &Memory{keys: map[string][]domain.Key{}} }
func (m *Memory) Keys(ctx context.Context, zone string) ([]domain.Key, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.Key{}, m.keys[zone]...), nil
}
func (m *Memory) Put(zone string, k domain.Key) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.keys[zone] = append(m.keys[zone], k)
}

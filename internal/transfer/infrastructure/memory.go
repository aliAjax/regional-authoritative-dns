package infrastructure

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/transfer/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	items []domain.Run
}

func New() *Memory { return &Memory{} }
func (m *Memory) Save(ctx context.Context, r domain.Run) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, x := range m.items {
		if x.ID == r.ID {
			m.items[i] = r
			return nil
		}
	}
	m.items = append(m.items, r)
	return nil
}
func (m *Memory) List(ctx context.Context) ([]domain.Run, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.Run{}, m.items...), nil
}

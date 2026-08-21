package infrastructure

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/worker/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	items []domain.Job
}

func (m *Memory) Save(ctx context.Context, j domain.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, x := range m.items {
		if x.ID == j.ID {
			m.items[i] = j
			return nil
		}
	}
	m.items = append(m.items, j)
	return nil
}
func (m *Memory) List(ctx context.Context) ([]domain.Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.Job{}, m.items...), nil
}

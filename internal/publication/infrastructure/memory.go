package infrastructure

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/publication/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	items []domain.Event
}

func (m *Memory) Save(ctx context.Context, e domain.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = append(m.items, e)
	return nil
}
func (m *Memory) List(ctx context.Context, z string) ([]domain.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var o []domain.Event
	for _, e := range m.items {
		if z == "" || e.ZoneID == z {
			o = append(o, e)
		}
	}
	return o, nil
}

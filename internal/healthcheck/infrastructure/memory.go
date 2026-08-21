package infrastructure

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/healthcheck/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	items []domain.Result
}

func (m *Memory) Put(ctx context.Context, r domain.Result) {
	if ctx.Err() != nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = append(m.items, r)
}
func (m *Memory) List(ctx context.Context) []domain.Result {
	if ctx.Err() != nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.Result{}, m.items...)
}

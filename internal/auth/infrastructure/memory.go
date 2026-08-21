package infrastructure

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/auth/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	items map[string][]domain.Permission
}

func New() *Memory { return &Memory{items: map[string][]domain.Permission{}} }
func (m *Memory) Permissions(ctx context.Context, s string) ([]domain.Permission, error) {
	// Respect a canceled context so a request that has already timed out does
	// not surface its permissions (and thus its subject) into audit records.
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.Permission{}, m.items[s]...), nil
}
func (m *Memory) Grant(s string, p domain.Permission) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[s] = append(m.items[s], p)
}

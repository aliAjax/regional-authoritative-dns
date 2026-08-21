package infrastructure

import (
	"context"
	"fmt"
	"github.com/example/regional-authoritative-dns/internal/record/domain"
	zoned "github.com/example/regional-authoritative-dns/internal/zone/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	zones map[string]zoned.Zone
}

func New() *Memory { return &Memory{zones: map[string]zoned.Zone{}} }
func (m *Memory) Create(ctx context.Context, z zoned.Zone) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.zones[z.ID]; ok {
		return fmt.Errorf("zone exists")
	}
	m.zones[z.ID] = z
	return nil
}
func (m *Memory) Get(ctx context.Context, id string) (zoned.Zone, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	z, ok := m.zones[id]
	if !ok {
		return z, fmt.Errorf("zone %s not found", id)
	}
	return z, nil
}
func (m *Memory) Save(ctx context.Context, z zoned.Zone) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.zones[z.ID]; !ok {
		return fmt.Errorf("zone %s not found", z.ID)
	}
	m.zones[z.ID] = z
	return nil
}
func (m *Memory) List(ctx context.Context) ([]zoned.Zone, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]zoned.Zone, 0, len(m.zones))
	for _, z := range m.zones {
		out = append(out, z)
	}
	return out, nil
}

type RecordMemory struct {
	mu      sync.RWMutex
	records map[string]domain.Set
}

func NewRecordMemory() *RecordMemory { return &RecordMemory{records: map[string]domain.Set{}} }
func (m *RecordMemory) Replace(ctx context.Context, id string, r domain.Set) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records[id] = append(domain.Set{}, r...)
	return nil
}
func (m *RecordMemory) List(ctx context.Context, id string) (domain.Set, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append(domain.Set{}, m.records[id]...), nil
}

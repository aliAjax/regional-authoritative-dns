package application

import (
	"context"
	"fmt"
	"github.com/example/regional-authoritative-dns/internal/platform"
	"github.com/example/regional-authoritative-dns/internal/record/domain"
	zoned "github.com/example/regional-authoritative-dns/internal/zone/domain"
	"sync"
	"time"
)

type Repository interface {
	Create(context.Context, zoned.Zone) error
	Get(context.Context, string) (zoned.Zone, error)
	Save(context.Context, zoned.Zone) error
	List(context.Context) ([]zoned.Zone, error)
}
type RecordRepository interface {
	Replace(context.Context, string, domain.Set) error
	List(context.Context, string) (domain.Set, error)
}
type Service struct {
	Zones   Repository
	Records RecordRepository
	mu      sync.Mutex
}

func New(z Repository, r RecordRepository) *Service { return &Service{Zones: z, Records: r} }
func (s *Service) Create(ctx context.Context, z zoned.Zone, records domain.Set) (zoned.Zone, error) {
	z.Name = zoned.NormalizeName(z.Name)
	for i := range records {
		records[i] = records[i].Normalize()
	}
	z.ID = platform.ID()
	z.Serial = 1
	z.Status = zoned.Draft
	z.DefaultTTL = defaultTTL(z.DefaultTTL)
	z.CreatedAt = time.Now().UTC()
	z.UpdatedAt = z.CreatedAt
	if e := z.Validate(); e != nil {
		return z, e
	}
	if e := records.Validate(); e != nil {
		return z, fmt.Errorf("records: %w", e)
	}
	if e := s.Zones.Create(ctx, z); e != nil {
		return z, fmt.Errorf("create zone: %w", e)
	}
	if e := s.Records.Replace(ctx, z.ID, records); e != nil {
		return z, fmt.Errorf("save records: %w", e)
	}
	return z, nil
}
func defaultTTL(v uint32) uint32 {
	if v == 0 {
		return 300
	}
	return v
}
func (s *Service) AddVersion(ctx context.Context, id string, records domain.Set, expected int64) (zoned.Zone, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	z, e := s.Zones.Get(ctx, id)
	if e != nil {
		return z, e
	}
	if expected > 0 && z.Version != expected {
		return z, fmt.Errorf("version conflict: expected %d got %d", expected, z.Version)
	}
	if !z.CanPublish() {
		return z, fmt.Errorf("zone status %s does not accept changes", z.Status)
	}
	if e := records.Validate(); e != nil {
		return z, e
	}
	z.Serial++
	z.Version++
	z.UpdatedAt = time.Now().UTC()
	if e := s.Records.Replace(ctx, id, records); e != nil {
		return z, e
	}
	if e := s.Zones.Save(ctx, z); e != nil {
		return z, e
	}
	return z, nil
}
func (s *Service) Publish(ctx context.Context, id string) (zoned.Zone, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	z, e := s.Zones.Get(ctx, id)
	if e != nil {
		return z, e
	}
	if _, e = s.Records.List(ctx, id); e != nil {
		return z, e
	}
	z.Status = zoned.Published
	z.Version++
	z.UpdatedAt = time.Now().UTC()
	if e = s.Zones.Save(ctx, z); e != nil {
		return z, e
	}
	return z, nil
}
func (s *Service) Freeze(ctx context.Context, id string) (zoned.Zone, error) {
	z, e := s.Zones.Get(ctx, id)
	if e != nil {
		return z, e
	}
	z.Status = zoned.Frozen
	z.UpdatedAt = time.Now().UTC()
	return z, s.Zones.Save(ctx, z)
}
func (s *Service) Rollback(ctx context.Context, id string, serial uint32) (zoned.Zone, error) {
	z, e := s.Zones.Get(ctx, id)
	if e != nil {
		return z, e
	}
	if serial == 0 {
		return z, fmt.Errorf("serial required")
	}
	z.Serial = serial
	z.Status = zoned.Published
	z.Version++
	z.UpdatedAt = time.Now().UTC()
	return z, s.Zones.Save(ctx, z)
}
func (s *Service) Get(ctx context.Context, id string) (zoned.Zone, domain.Set, error) {
	z, e := s.Zones.Get(ctx, id)
	if e != nil {
		return z, nil, e
	}
	r, e := s.Records.List(ctx, id)
	return z, r, e
}
func (s *Service) List(ctx context.Context) ([]zoned.Zone, error) { return s.Zones.List(ctx) }

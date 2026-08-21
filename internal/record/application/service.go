package application

import (
	"context"
	"fmt"
	"github.com/example/regional-authoritative-dns/internal/record/domain"
	"sort"
	"strings"
)

type Repository interface {
	Replace(context.Context, string, domain.Set) error
	List(context.Context, string) (domain.Set, error)
}
type Service struct{ Repo Repository }

func New(r Repository) *Service { return &Service{Repo: r} }
func (s *Service) Upsert(ctx context.Context, zone string, records domain.Set) error {
	if zone == "" {
		return fmt.Errorf("zone required")
	}
	normalized := make(domain.Set, 0, len(records))
	seen := map[string]bool{}
	for _, r := range records {
		if seen[r.Key()] {
			continue
		}
		seen[r.Key()] = true
		r = r.Normalize()
		if e := r.Validate(); e != nil {
			return fmt.Errorf("%s: %w", r.Name, e)
		}
		normalized = append(normalized, r)
	}
	sort.Slice(normalized, func(i, j int) bool { return normalized[i].Key() < normalized[j].Key() })
	return s.Repo.Replace(ctx, zone, normalized)
}
func (s *Service) Filter(ctx context.Context, zone, name string, typ domain.Type, view string) (domain.Set, error) {
	rs, e := s.Repo.List(ctx, zone)
	if e != nil {
		return nil, e
	}
	var out domain.Set
	for _, r := range rs {
		if name != "" && strings.ToLower(r.Name) != strings.ToLower(name) {
			continue
		}
		if typ != "" && r.Type != typ {
			continue
		}
		if view != "" && r.View != view {
			continue
		}
		out = append(out, r)
	}
	return out, nil
}
func (s *Service) Count(ctx context.Context, zone string) (int, error) {
	rs, e := s.Repo.List(ctx, zone)
	return len(rs), e
}

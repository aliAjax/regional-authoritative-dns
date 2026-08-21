package application

import (
	"context"
	"fmt"
	"github.com/example/regional-authoritative-dns/internal/auth/domain"
)

type Repository interface {
	Permissions(context.Context, string) ([]domain.Permission, error)
}
type Service struct{ Repo Repository }

func New(r Repository) *Service { return &Service{Repo: r} }
func (s *Service) Authorize(ctx context.Context, subject, zone, action string) error {
	ps, e := s.Repo.Permissions(ctx, subject)
	if e != nil {
		return fmt.Errorf("load permissions: %v", e)
	}
	for _, p := range ps {
		if p.Matches(zone, action) {
			return nil
		}
	}
	return fmt.Errorf("permission denied")
}

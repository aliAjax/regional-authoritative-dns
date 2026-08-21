package application

import (
	"context"
	"fmt"
	"github.com/example/regional-authoritative-dns/internal/platform"
	"github.com/example/regional-authoritative-dns/internal/transfer/domain"
	"sync"
	"time"
)

type Repository interface {
	Save(context.Context, domain.Run) error
	List(context.Context) ([]domain.Run, error)
}
type Service struct {
	Repo Repository
	mu   sync.Mutex
}

func New(r Repository) *Service { return &Service{Repo: r} }
func (s *Service) Start(ctx context.Context, zone, target string, mode domain.Mode, from, to uint32) (domain.Run, error) {
	if zone == "" || target == "" {
		return domain.Run{}, fmt.Errorf("zone and target required")
	}
	r := domain.Run{ID: platform.ID(), ZoneID: zone, Target: target, Mode: mode, FromSerial: from, ToSerial: to, Status: "running", StartedAt: time.Now().UTC()}
	if e := s.Repo.Save(ctx, r); e != nil {
		return r, e
	}
	return r, nil
}
func (s *Service) Complete(ctx context.Context, r domain.Run, e error) error {
	if e != nil {
		r.Status = "failed"
		r.Error = e.Error()
		r.Attempts++
	} else {
		r.Status = "completed"
	}
	r.FinishedAt = time.Now().UTC()
	return s.Repo.Save(ctx, r)
}

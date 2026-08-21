package application

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/platform"
	"github.com/example/regional-authoritative-dns/internal/worker/domain"
	"log/slog"
	"sync"
	"time"
)

type Repository interface {
	Save(context.Context, domain.Job) error
	List(context.Context) ([]domain.Job, error)
}
type Handler func(context.Context, domain.Job) error
type Runner struct {
	Repo     Repository
	Handlers map[string]Handler
	Interval time.Duration
	logger   *slog.Logger
	mu       sync.Mutex
}

func New(r Repository, l *slog.Logger) *Runner {
	return &Runner{Repo: r, Handlers: map[string]Handler{}, Interval: 30 * time.Second, logger: l}
}
func (s *Runner) Enqueue(ctx context.Context, kind, zone string) domain.Job {
	j := domain.Job{ID: platform.ID(), Kind: kind, ZoneID: zone, Status: "queued", NextRun: time.Now().UTC()}
	_ = s.Repo.Save(ctx, j)
	return j
}
func (s *Runner) Run(ctx context.Context) {
	t := time.NewTicker(s.Interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			s.once(ctx)
		}
	}
}
func (s *Runner) once(ctx context.Context) {
	jobs, _ := s.Repo.List(ctx)
	for _, j := range jobs {
		h := s.Handlers[j.Kind]
		if h == nil || j.Status == "completed" || j.Status == "failed" {
			continue
		}
		j.Status = "running"
		if e := h(ctx, j); e != nil {
			j.Status = "failed"
			j.Error = e.Error()
			j.Attempts++
		} else {
			j.Status = "completed"
		}
		_ = s.Repo.Save(ctx, j)
	}
}

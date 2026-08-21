package application

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/platform"
	"github.com/example/regional-authoritative-dns/internal/publication/domain"
	"time"
)

type Repository interface {
	Save(context.Context, domain.Event) error
	List(context.Context, string) ([]domain.Event, error)
}
type Service struct{ Repo Repository }

func New(r Repository) *Service { return &Service{Repo: r} }
func (s *Service) Record(ctx context.Context, zone, action, actor, reason, request, summary string) error {
	if e := ctx.Err(); e != nil {
		// A canceled request must not append to the audit log, since the event
		// it describes never completed and persisting it would mislead readers.
		return e
	}
	return s.Repo.Save(ctx, domain.Event{ID: platform.ID(), ZoneID: zone, Action: action, Actor: actor, Reason: reason, RequestID: request, Summary: summary, CreatedAt: time.Now().UTC()})
}

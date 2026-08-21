package application

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/publication/domain"
	"testing"
)

type cancelRepo struct{}

func (cancelRepo) Save(context.Context, domain.Event) error             { return nil }
func (cancelRepo) List(context.Context, string) ([]domain.Event, error) { return nil, nil }
func TestRecordRejectsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := New(cancelRepo{}).Record(ctx, "z", "publish", "a", "r", "q", "s"); e == nil {
		t.Fatal("canceled publication recorded")
	}
}

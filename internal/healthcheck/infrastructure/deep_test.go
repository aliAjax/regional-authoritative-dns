package infrastructure

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/healthcheck/domain"
	"testing"
)

func TestListRejectsCanceledContext(t *testing.T) {
	m := &Memory{}
	m.Put(context.Background(), domain.Result{Name: "dns"})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if got := m.List(ctx); len(got) != 0 {
		t.Fatalf("canceled list returned %#v", got)
	}
}

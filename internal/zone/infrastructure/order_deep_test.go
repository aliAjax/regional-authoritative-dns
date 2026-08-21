package infrastructure

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/zone/domain"
	"testing"
)

func TestListReturnsStableZoneOrder(t *testing.T) {
	m := New()
	_ = m.Create(context.Background(), domain.Zone{ID: "b", Name: "b.example."})
	_ = m.Create(context.Background(), domain.Zone{ID: "a", Name: "a.example."})
	got, _ := m.List(context.Background())
	if len(got) != 2 || got[0].ID != "a" {
		t.Fatalf("zone order unstable: %#v", got)
	}
}

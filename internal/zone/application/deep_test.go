package application

import (
	"context"
	record "github.com/example/regional-authoritative-dns/internal/record/domain"
	"github.com/example/regional-authoritative-dns/internal/zone/domain"
	"github.com/example/regional-authoritative-dns/internal/zone/infrastructure"
	"testing"
)

func TestFrozenZoneRejectsVersionChangesAndPublishing(t *testing.T) {
	s := New(infrastructure.New(), infrastructure.NewRecordMemory())
	z, e := s.Create(context.Background(), domain.Zone{Name: "frozen.example"}, record.Set{{Name: "www.frozen.example.", Type: record.A, TTL: 60, Data: "192.0.2.3"}})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Freeze(context.Background(), z.ID); e != nil {
		t.Fatal(e)
	}
	if _, e = s.AddVersion(context.Background(), z.ID, record.Set{{Name: "www.frozen.example.", Type: record.A, TTL: 60, Data: "192.0.2.4"}}, 0); e == nil {
		t.Fatal("frozen version accepted")
	}
	if _, e = s.Publish(context.Background(), z.ID); e == nil {
		t.Fatal("frozen publish accepted")
	}
}

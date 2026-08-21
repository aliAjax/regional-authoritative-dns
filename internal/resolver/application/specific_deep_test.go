package application

import (
	"context"
	record "github.com/example/regional-authoritative-dns/internal/record/domain"
	resdomain "github.com/example/regional-authoritative-dns/internal/resolver/domain"
	zone "github.com/example/regional-authoritative-dns/internal/zone/application"
	zoned "github.com/example/regional-authoritative-dns/internal/zone/domain"
	zoneinfra "github.com/example/regional-authoritative-dns/internal/zone/infrastructure"
	"testing"
)

func TestResolveNameUsesMostSpecificAuthoritativeZone(t *testing.T) {
	zr, rr := zoneinfra.New(), zoneinfra.NewRecordMemory()
	zs := zone.New(zr, rr)
	ctx := context.Background()
	_, e := zs.Create(ctx, zoned.Zone{Name: "example.com"}, record.Set{{Name: "api.example.com.", Type: record.A, TTL: 60, Data: "192.0.2.1"}})
	if e != nil {
		t.Fatal(e)
	}
	_, e = zs.Create(ctx, zoned.Zone{Name: "sub.example.com"}, record.Set{{Name: "api.sub.example.com.", Type: record.A, TTL: 60, Data: "192.0.2.2"}})
	if e != nil {
		t.Fatal(e)
	}
	r, e := New(zs).ResolveName(ctx, resdomain.Query{Name: "api.sub.example.com.", Type: record.A})
	if e != nil || len(r.Answers) != 1 || r.Answers[0].Data != "192.0.2.2" {
		t.Fatalf("wrong zone answer: %#v %v", r, e)
	}
}

package application

import (
	"context"
	record "github.com/example/regional-authoritative-dns/internal/record/domain"
	"testing"
)

type deepRepo struct{ saved record.Set }

func (r *deepRepo) Replace(_ context.Context, _ string, rs record.Set) error {
	r.saved = append(record.Set{}, rs...)
	return nil
}
func (r *deepRepo) List(context.Context, string) (record.Set, error) {
	return append(record.Set{}, r.saved...), nil
}

func TestUpsertDeduplicatesNormalizedRecords(t *testing.T) {
	r := &deepRepo{}
	if err := New(r).Upsert(context.Background(), "zone", record.Set{{Name: "WWW.Example", Type: record.A, TTL: 60, Data: "192.0.2.1"}, {Name: "www.example.", Type: record.A, TTL: 60, Data: "192.0.2.1"}}); err != nil {
		t.Fatal(err)
	}
	if len(r.saved) != 1 {
		t.Fatalf("normalized duplicate was stored: %#v", r.saved)
	}
}

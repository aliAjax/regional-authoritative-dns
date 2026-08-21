package application

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/dnssec/domain"
	record "github.com/example/regional-authoritative-dns/internal/record/domain"
	"testing"
)

type emptyKeys struct{}

func (emptyKeys) Keys(context.Context, string) ([]domain.Key, error) { return nil, nil }
func TestVerifyRejectsEmptyRecordSet(t *testing.T) {
	if r, e := New(nil).Verify(context.Background(), "example.", "www.example.", record.Set{}); e == nil && r.Valid {
		t.Fatal("empty set verified")
	}
}

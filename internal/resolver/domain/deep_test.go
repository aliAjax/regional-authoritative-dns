package domain

import (
	record "github.com/example/regional-authoritative-dns/internal/record/domain"
	"testing"
)

func TestQueryNormalizeAppliesDNSDefaults(t *testing.T) {
	q := Query{Name: "WWW.Example", Type: record.A}
	q = q.Normalize()
	if q.Name != "www.example." || q.View != "public" {
		t.Fatalf("query defaults missing: %#v", q)
	}
}

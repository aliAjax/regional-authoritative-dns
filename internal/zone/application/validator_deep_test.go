package application

import (
	record "github.com/example/regional-authoritative-dns/internal/record/domain"
	"github.com/example/regional-authoritative-dns/internal/zone/domain"
	"testing"
)

func TestValidateDetectsCNAMEConflict(t *testing.T) {
	v := Validate(domain.Zone{Name: "x.example.", Serial: 1, DefaultTTL: 60}, record.Set{{Name: "a.x.example.", Type: record.CNAME, TTL: 60, Data: "target.example."}, {Name: "a.x.example.", Type: record.A, TTL: 60, Data: "192.0.2.1"}})
	if v.Valid() {
		t.Fatal("CNAME conflict was accepted")
	}
}

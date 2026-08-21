package adapter

import (
	record "github.com/example/regional-authoritative-dns/internal/record/domain"
	"testing"
)

func TestDigestIsStableAcrossRecordOrderAndNameCase(t *testing.T) {
	a := record.Set{{Name: "WWW.Example.", Type: record.A, TTL: 60, Data: "192.0.2.5"}, {Name: "mail.example.", Type: record.MX, TTL: 60, Data: "10 mx.example."}}
	b := record.Set{{Name: "mail.example.", Type: record.MX, TTL: 60, Data: "10 mx.example."}, {Name: "www.example.", Type: record.A, TTL: 60, Data: "192.0.2.5"}}
	if Digest(a) != Digest(b) {
		t.Fatal("equivalent records have different digest")
	}
}

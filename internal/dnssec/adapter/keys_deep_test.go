package adapter

import (
	"github.com/example/regional-authoritative-dns/internal/dnssec/domain"
	"testing"
)

func TestValidateKeyRequiresActiveKey(t *testing.T) {
	if err := ValidateKey(domain.Key{ID: "k", Public: "p", Role: domain.KSK, Active: false}); err == nil {
		t.Fatal("inactive key accepted")
	}
}

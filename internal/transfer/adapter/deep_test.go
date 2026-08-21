package adapter

import (
	"github.com/example/regional-authoritative-dns/internal/transfer/domain"
	"testing"
)

func TestValidateModeRejectsUnknownMode(t *testing.T) {
	if ValidateMode(domain.Mode("bad")) == nil {
		t.Fatal("unknown transfer mode accepted")
	}
}

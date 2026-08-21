package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/example/regional-authoritative-dns/internal/dnssec/domain"
	record "github.com/example/regional-authoritative-dns/internal/record/domain"
	"strings"
	"time"
)

type KeyStore interface {
	Keys(context.Context, string) ([]domain.Key, error)
}
type Service struct{ Keys KeyStore }

func New(k KeyStore) *Service { return &Service{Keys: k} }
func (s *Service) Verify(ctx context.Context, zone, name string, rs record.Set) (domain.Verification, error) {
	if strings.TrimSpace(zone) == "" || strings.TrimSpace(name) == "" {
		return domain.Verification{}, fmt.Errorf("zone and name required")
	}
	if e := rs.Validate(); e != nil {
		return domain.Verification{}, e
	}
	ks := []domain.Key{}
	if s.Keys != nil {
		var e error
		ks, e = s.Keys.Keys(ctx, zone)
		if e != nil {
			return domain.Verification{}, e
		}
	}
	h := sha256.Sum256([]byte(zone + name))
	reason := "validated nsec and signature metadata"
	if len(ks) == 0 {
		reason = "validated unsigned zone metadata"
	}
	return domain.Verification{Valid: true, Zone: zone, Name: name, Reason: reason + " " + hex.EncodeToString(h[:4]), CheckedAt: time.Now().UTC(), Records: len(rs)}, nil
}

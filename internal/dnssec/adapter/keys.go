package adapter

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/example/regional-authoritative-dns/internal/dnssec/domain"
)

func GenerateKey(role domain.KeyRole, algorithm string) (domain.Key, error) {
	b := make([]byte, 32)
	if _, e := rand.Read(b); e != nil {
		return domain.Key{}, e
	}
	return domain.Key{ID: base64.RawURLEncoding.EncodeToString(b[:8]), Role: role, Algorithm: algorithm, Public: base64.RawURLEncoding.EncodeToString(b), Active: true}, nil
}
func ValidateKey(k domain.Key) error {
	if k.ID == "" || k.Public == "" {
		return fmt.Errorf("key identity required")
	}
	if k.Role != domain.KSK && k.Role != domain.ZSK {
		return fmt.Errorf("unknown role")
	}
	if !k.Active {
		return fmt.Errorf("key is not active")
	}
	return nil
}

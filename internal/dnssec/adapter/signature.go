package adapter

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/example/regional-authoritative-dns/internal/record/domain"
	"strings"
	"time"
)

type Signer interface {
	Sign([]byte) ([]byte, error)
	Algorithm() string
}

func Canonical(r domain.Record) string {
	return strings.ToLower(r.Name) + " " + strings.ToUpper(string(r.Type)) + " " + r.Data
}
func Digest(rs domain.Set) string {
	h := sha256.New()
	for _, r := range rs.Sort() {
		_, _ = h.Write([]byte(Canonical(r)))
		_, _ = h.Write([]byte{'\n'})
	}
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
func BuildRRSIG(r domain.Record, keytag uint16, alg uint8, expires time.Time) domain.Record {
	return domain.Record{Name: r.Name, Type: domain.RRSIG, TTL: r.TTL, Data: fmt.Sprintf("%s %d %s", r.Type, keytag, expires.UTC().Format("20060102150405"))}
}
func BuildNSEC(name, next string, types []domain.Type, ttl uint32) domain.Record {
	return domain.Record{Name: name, Type: domain.NSEC, TTL: ttl, Data: next + " " + strings.Join(func() []string {
		out := []string{}
		for _, t := range types {
			out = append(out, string(t))
		}
		return out
	}(), " ")}
}

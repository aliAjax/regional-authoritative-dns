package application

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/example/regional-authoritative-dns/internal/record/domain"
	"sort"
	"strings"
)

func Digest(rs domain.Set) string {
	items := make([]string, 0, len(rs))
	for _, r := range rs {
		items = append(items, r.Key()+"|"+string(r.Type)+"|"+r.Data)
	}
	sort.Strings(items)
	h := sha256.Sum256([]byte(strings.Join(items, "\n")))
	return hex.EncodeToString(h[:])
}
func Summary(rs domain.Set) map[string]int {
	out := map[string]int{}
	for _, r := range rs {
		out[string(r.Type)]++
	}
	return out
}

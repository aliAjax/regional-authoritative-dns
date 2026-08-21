package adapter

import (
	"github.com/example/regional-authoritative-dns/internal/record/domain"
	"net"
	"sort"
	"strings"
)

type ViewPolicy struct {
	Name     string
	Networks []*net.IPNet
	Priority int
	Enabled  bool
}

func (p ViewPolicy) Matches(ip string) bool {
	if !p.Enabled {
		return false
	}
	if len(p.Networks) == 0 {
		return true
	}
	addr := net.ParseIP(ip)
	for _, n := range p.Networks {
		if addr != nil && n.Contains(addr) {
			return true
		}
	}
	return false
}
func SelectView(p []ViewPolicy, ip string) string {
	// Sort a private copy so the caller's slice is never reordered. The
	// original sort mutated the shared input in place, which both surprised
	// callers and raced when SelectView was invoked concurrently.
	sorted := make([]ViewPolicy, len(p))
	copy(sorted, p)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Priority > sorted[j].Priority })
	for _, x := range sorted {
		if x.Matches(ip) {
			return x.Name
		}
	}
	return "public"
}
func SelectWeighted(rs domain.Set) []domain.Record {
	var healthy domain.Set
	for _, r := range rs {
		if r.Healthy {
			healthy = append(healthy, r)
		}
	}
	sort.SliceStable(healthy, func(i, j int) bool {
		if healthy[i].Priority == healthy[j].Priority {
			return healthy[i].Weight > healthy[j].Weight
		}
		return healthy[i].Priority < healthy[j].Priority
	})
	return healthy
}
func CanonicalName(s string) string {
	return strings.ToLower(strings.TrimSuffix(strings.TrimSpace(s), ".")) + "."
}

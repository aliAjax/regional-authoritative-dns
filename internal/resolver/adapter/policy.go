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
	sort.Slice(p, func(i, j int) bool { return p[i].Priority > p[j].Priority })
	for _, x := range p {
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

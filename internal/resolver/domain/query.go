package domain

import (
	"strings"

	"github.com/example/regional-authoritative-dns/internal/record/domain"
)

type Query struct {
	Name     string
	Type     domain.Type
	View     string
	ClientIP string
	TCP      bool
}
type Answer struct {
	Name string
	Type domain.Type
	TTL  uint32
	Data string
}
type Result struct {
	Answers   []Answer
	Authority []Answer
	NXDomain  bool
	NoData    bool
	ServFail  bool
	Truncated bool
	Serial    uint32
}

// Normalize applies the DNS query defaults: the name is lowercased and made
// fully-qualified with a trailing dot, and an empty view falls back to the
// public view. The remaining fields are returned untouched.
func (q Query) Normalize() Query {
	q.Name = strings.ToLower(strings.TrimSpace(q.Name))
	if q.Name != "" && !strings.HasSuffix(q.Name, ".") {
		q.Name += "."
	}
	if q.View == "" {
		q.View = "public"
	}
	return q
}

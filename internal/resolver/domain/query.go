package domain

import "github.com/example/regional-authoritative-dns/internal/record/domain"

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

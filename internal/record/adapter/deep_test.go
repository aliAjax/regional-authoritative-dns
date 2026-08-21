package adapter

import (
	record "github.com/example/regional-authoritative-dns/internal/record/domain"
	"testing"
)

func TestParseLineCanonicalizesType(t *testing.T) {
	r, err := ParseLine("www.example. 60 IN a 192.0.2.1")
	if err != nil || r.Type != record.A {
		t.Fatalf("lowercase record type was not canonicalized: %#v %v", r, err)
	}
}

package adapter

import (
	"github.com/example/regional-authoritative-dns/internal/healthcheck/domain"
	"testing"
	"time"
)

func TestExpiredTreatsFutureRefreshAsNotExpired(t *testing.T) {
	s := New()
	s.Put(domain.Result{Name: "dns", Status: domain.Healthy, CheckedAt: time.Now().Add(time.Second)})
	if s.Expired("dns", 100*time.Millisecond) {
		t.Fatal("future health result was expired")
	}
}

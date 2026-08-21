package application

import (
	"context"
	"fmt"
	"github.com/example/regional-authoritative-dns/internal/healthcheck/domain"
	"net"
	"time"
)

type Service struct{ Timeout time.Duration }

func New(d time.Duration) *Service {
	if d <= 0 {
		d = 2 * time.Second
	}
	return &Service{Timeout: d}
}
func (s *Service) Check(ctx context.Context, name, address string) domain.Result {
	start := time.Now()
	r := domain.Result{Name: name, CheckedAt: start}
	c, e := net.DialTimeout("tcp", address, s.Timeout)
	if e != nil {
		r.Status = domain.Unhealthy
		r.Error = fmt.Sprintf("connect: %v", e)
	} else {
		r.Status = domain.Healthy
		_ = c.Close()
	}
	r.Latency = time.Since(start)
	return r
}

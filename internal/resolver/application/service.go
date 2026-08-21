package application

import (
	"context"
	"fmt"
	resdomain "github.com/example/regional-authoritative-dns/internal/resolver/domain"
	"github.com/example/regional-authoritative-dns/internal/zone/application"
	"math/rand"
	"strings"
)

type Service struct{ Zones *application.Service }

func New(z *application.Service) *Service { return &Service{Zones: z} }
func (s *Service) Resolve(ctx context.Context, zoneName string, q resdomain.Query) (resdomain.Result, error) {
	z, records, e := s.Zones.Get(ctx, zoneName)
	if e != nil {
		zones, le := s.Zones.List(ctx)
		if le != nil {
			return resdomain.Result{ServFail: true}, e
		}
		foundZone := false
		for _, candidate := range zones {
			if candidate.Name == zoneName {
				z, records, e = s.Zones.Get(ctx, candidate.ID)
				foundZone = true
				break
			}
		}
		if !foundZone || e != nil {
			return resdomain.Result{ServFail: true}, e
		}
	}
	q.Name = strings.ToLower(q.Name)
	if !strings.HasSuffix(q.Name, ".") {
		q.Name += "."
	}
	if q.View == "" {
		q.View = "public"
	}
	found := records.Find(q.Name, q.Type, q.View)
	if len(found) == 0 {
		for _, r := range records {
			if r.Type == q.Type && r.View == q.View && strings.HasPrefix(r.Name, "*.") && strings.HasSuffix(q.Name, strings.TrimPrefix(r.Name, "*")) {
				found = append(found, r)
				break
			}
		}
	}
	out := resdomain.Result{Serial: z.Serial}
	if len(found) == 0 {
		nameExists := false
		for _, r := range records {
			if r.Name == q.Name {
				nameExists = true
				break
			}
		}
		if nameExists {
			out.NoData = true
		} else {
			out.NXDomain = true
		}
		return out, nil
	}
	for _, r := range found {
		if !r.Healthy {
			continue
		}
		out.Answers = append(out.Answers, resdomain.Answer{Name: r.Name, Type: r.Type, TTL: r.TTL, Data: r.Data})
	}
	if len(out.Answers) == 0 {
		out.ServFail = true
		return out, nil
	}
	if len(out.Answers) > 1 {
		rand.Shuffle(len(out.Answers), func(i, j int) { out.Answers[i], out.Answers[j] = out.Answers[j], out.Answers[i] })
	}
	return out, nil
}

func (s *Service) ResolveName(ctx context.Context, q resdomain.Query) (resdomain.Result, error) {
	for _, z := range func() []struct{ ID, Name string } {
		zones, _ := s.Zones.List(ctx)
		out := make([]struct{ ID, Name string }, 0, len(zones))
		for _, z := range zones {
			out = append(out, struct{ ID, Name string }{z.ID, z.Name})
		}
		return out
	}() {
		if strings.HasSuffix(strings.ToLower(q.Name), strings.TrimPrefix(strings.ToLower(z.Name), ".")) {
			if r, e := s.Resolve(ctx, z.ID, q); e == nil && !r.ServFail {
				return r, nil
			}
		}
	}
	return resdomain.Result{ServFail: true}, fmt.Errorf("authoritative zone not found for %s", q.Name)
}
func (s *Service) ValidateZone(ctx context.Context, id string) error {
	_, r, e := s.Zones.Get(ctx, id)
	if e != nil {
		return e
	}
	if e = r.Validate(); e != nil {
		return fmt.Errorf("validate zone: %w", e)
	}
	return nil
}

package application

import (
	"fmt"
	"github.com/example/regional-authoritative-dns/internal/record/domain"
	zoned "github.com/example/regional-authoritative-dns/internal/zone/domain"
	"strings"
)

type Validation struct {
	Errors   []string
	Warnings []string
}

func Validate(z zoned.Zone, rs domain.Set) Validation {
	v := Validation{}
	if e := z.Validate(); e != nil {
		v.Errors = append(v.Errors, e.Error())
	}
	if e := rs.Validate(); e != nil {
		v.Errors = append(v.Errors, e.Error())
	}
	names := map[string]map[domain.Type]int{}
	for _, r := range rs {
		if names[r.Name] == nil {
			names[r.Name] = map[domain.Type]int{}
		}
		names[r.Name][r.Type]++
		if r.Type == domain.CNAME && len(names[r.Name]) > 2 {
			v.Errors = append(v.Errors, fmt.Sprintf("CNAME conflict at %s", r.Name))
		}
		if r.TTL < 30 {
			v.Warnings = append(v.Warnings, fmt.Sprintf("low TTL at %s", r.Name))
		}
		if strings.Contains(r.Data, " ") && r.Type == domain.A {
			v.Errors = append(v.Errors, "A record has spaces")
		}
	}
	return v
}
func (v Validation) Valid() bool { return len(v.Errors) == 0 }
func (v Validation) Error() error {
	if len(v.Errors) == 0 {
		return nil
	}
	return fmt.Errorf("zone validation failed: %s", strings.Join(v.Errors, "; "))
}

package domain

import (
	"fmt"
	"strings"
	"time"
)

type Status string

const (
	Draft     Status = "draft"
	Published Status = "published"
	Frozen    Status = "frozen"
	Archived  Status = "archived"
)

type Zone struct {
	ID, Name, Description string
	Serial                uint32
	Status                Status
	DefaultTTL            uint32
	Views                 []string
	CreatedAt, UpdatedAt  time.Time
	Version               int64
}

func NormalizeName(n string) string {
	n = strings.TrimSpace(strings.ToLower(n))
	if n == "" {
		return "."
	}
	if !strings.HasSuffix(n, ".") {
		n += "."
	}
	return n
}
func (z Zone) Validate() error {
	z.Name = NormalizeName(z.Name)
	if z.Name == "." {
		return fmt.Errorf("zone name required")
	}
	if z.DefaultTTL == 0 || z.DefaultTTL > 86400*30 {
		return fmt.Errorf("invalid default ttl")
	}
	if z.Serial == 0 {
		return fmt.Errorf("serial must be positive")
	}
	return nil
}
func (z Zone) CanPublish() bool { return z.Status != Archived }

package domain

import "time"

type Status string

const (
	Healthy   Status = "healthy"
	Unhealthy Status = "unhealthy"
)

type Result struct {
	Name      string
	Status    Status
	Latency   time.Duration
	Error     string
	CheckedAt time.Time
}

// IsHealthy reports whether the result represents a healthy target. A result is
// healthy only when the status is Healthy and no error was recorded, so a result
// that carries a stale Healthy status alongside an error is never accepted.
func (r Result) IsHealthy() bool {
	return r.Status == Healthy && r.Error == ""
}

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

package domain

import "time"

type Mode string

const (
	AXFR Mode = "AXFR"
	IXFR Mode = "IXFR"
)

type Run struct {
	ID, ZoneID, Target    string
	Mode                  Mode
	FromSerial, ToSerial  uint32
	Status                string
	Attempts              int
	Error                 string
	StartedAt, FinishedAt time.Time
}

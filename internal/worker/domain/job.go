package domain

import "time"

type Job struct {
	ID, Kind, ZoneID, Status string
	Attempts                 int
	NextRun                  time.Time
	Error                    string
}

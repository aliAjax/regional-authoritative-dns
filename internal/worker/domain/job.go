package domain

import "time"

type Job struct {
	ID, Kind, ZoneID, Status string
	Attempts                 int
	NextRun                  time.Time
	Error                    string
}

func (j Job) Runnable(now time.Time) bool {
	return j.Status != "completed" && !j.NextRun.After(now)
}

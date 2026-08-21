package domain

import "time"

type Event struct {
	ID, ZoneID, Action, Actor, Reason, RequestID, Summary string
	CreatedAt                                             time.Time
}

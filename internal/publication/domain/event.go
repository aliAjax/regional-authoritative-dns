package domain

import "time"

type Event struct {
	ID, ZoneID, Action, Actor, Reason, RequestID, Summary string
	CreatedAt                                             time.Time
}

// Ready reports whether the event carries the minimum identifying fields to be
// safely published: a generated identity and the action it records. Events
// missing either are incomplete and must not be appended to the audit log.
func (e Event) Ready() bool {
	return e.ID != "" && e.Action != ""
}

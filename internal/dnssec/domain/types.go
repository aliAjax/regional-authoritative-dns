package domain

import "time"

type KeyRole string

const (
	KSK KeyRole = "ksk"
	ZSK KeyRole = "zsk"
)

type Key struct {
	ID, Algorithm        string
	Role                 KeyRole
	Public               string
	Active               bool
	CreatedAt, ExpiresAt time.Time
}
type Verification struct {
	Valid              bool
	Zone, Name, Reason string
	CheckedAt          time.Time
	Records            int
}

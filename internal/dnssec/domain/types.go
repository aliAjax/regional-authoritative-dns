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

func (k Key) ActiveAt(now time.Time) bool {
	if !k.Active {
		return false
	}
	if !k.CreatedAt.IsZero() && now.Before(k.CreatedAt) {
		return false
	}
	if !k.ExpiresAt.IsZero() && now.After(k.ExpiresAt) {
		return false
	}
	return true
}

type Verification struct {
	Valid              bool
	Zone, Name, Reason string
	CheckedAt          time.Time
	Records            int
}

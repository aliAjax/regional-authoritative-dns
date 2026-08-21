package domain

import (
	"testing"
	"time"
)

func TestCanPublishRejectsFrozenAndArchived(t *testing.T) {
	for _, s := range []Status{Frozen, Archived} {
		if (Zone{Status: s}).CanPublish() {
			t.Fatalf("status %s can publish", s)
		}
	}
}
func TestVersionPublishedNeedsTimestamp(t *testing.T) {
	v := Version{Status: Published}
	if v.Published() {
		t.Fatal("version without publication time is published")
	}
	now := time.Now()
	v.PublishedAt = &now
	if !v.Published() {
		t.Fatal("timestamped version not published")
	}
}

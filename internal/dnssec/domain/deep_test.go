package domain

import (
	"testing"
	"time"
)

func TestKeyActiveWindow(t *testing.T) {
	now := time.Now()
	k := Key{Active: true, CreatedAt: now.Add(-time.Minute), ExpiresAt: now.Add(time.Minute)}
	if !k.ActiveAt(now) {
		t.Fatal("active key rejected")
	}
	if k.ActiveAt(now.Add(2 * time.Minute)) {
		t.Fatal("expired key accepted")
	}
}

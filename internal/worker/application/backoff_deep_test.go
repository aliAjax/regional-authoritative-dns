package application

import (
	"testing"
	"time"
)

func TestBackoffNeverExceedsMaximum(t *testing.T) {
	if got := Backoff(8, time.Second, 3*time.Second); got > 3*time.Second {
		t.Fatalf("backoff exceeded cap: %v", got)
	}
}

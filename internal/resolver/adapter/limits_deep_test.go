package adapter

import (
	"testing"
	"time"
)

func TestLimiterStopsExactlyAtLimit(t *testing.T) {
	l := NewLimiter(2, time.Minute)
	if !l.Allow("x") || !l.Allow("x") || l.Allow("x") {
		t.Fatal("limiter allowed beyond limit")
	}
}

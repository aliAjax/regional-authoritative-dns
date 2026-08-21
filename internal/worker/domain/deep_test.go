package domain

import (
	"testing"
	"time"
)

func TestJobRunnableRequiresDueTime(t *testing.T) {
	j := Job{Status: "queued", NextRun: time.Now().Add(time.Hour)}
	if j.Runnable(time.Now()) {
		t.Fatal("future job marked runnable")
	}
}

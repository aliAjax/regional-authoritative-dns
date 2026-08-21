package application

import (
	"context"
	"testing"
	"time"
)

func TestCheckHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := New(time.Second).Check(ctx, "dns", "192.0.2.250:9")
	if r.Status != "unhealthy" || r.Error == "" {
		t.Fatalf("canceled check returned %#v", r)
	}
}

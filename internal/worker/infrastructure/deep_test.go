package infrastructure

import (
	"context"
	worker "github.com/example/regional-authoritative-dns/internal/worker/domain"
	"testing"
)

func TestSaveRejectsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := new(Memory).Save(ctx, worker.Job{ID: "j"}); err == nil {
		t.Fatal("canceled save succeeded")
	}
}

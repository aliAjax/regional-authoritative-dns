package infrastructure

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/publication/domain"
	"testing"
)

func TestSaveRejectsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := new(Memory).Save(ctx, domain.Event{ID: "e"}); e == nil {
		t.Fatal("canceled event saved")
	}
}

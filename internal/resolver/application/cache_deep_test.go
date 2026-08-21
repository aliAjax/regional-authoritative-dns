package application

import (
	"context"
	resdomain "github.com/example/regional-authoritative-dns/internal/resolver/domain"
	"testing"
	"time"
)

func TestCacheGetHonorsCanceledContext(t *testing.T) {
	c := NewCache(time.Minute)
	c.Put(context.Background(), "k", resdomain.Result{Serial: 7})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, ok := c.Get(ctx, "k"); ok {
		t.Fatal("canceled cache hit")
	}
}
func TestCachePutIgnoresCanceledContext(t *testing.T) {
	c := NewCache(time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c.Put(ctx, "k", resdomain.Result{Serial: 7})
	if _, ok := c.Get(context.Background(), "k"); ok {
		t.Fatal("canceled cache write stored value")
	}
}

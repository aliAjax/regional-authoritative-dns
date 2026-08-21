package application

import (
	"context"
	resdomain "github.com/example/regional-authoritative-dns/internal/resolver/domain"
	"sync"
	"time"
)

type Cache struct {
	mu    sync.RWMutex
	ttl   time.Duration
	items map[string]entry
}
type entry struct {
	result  resdomain.Result
	expires time.Time
}

func NewCache(ttl time.Duration) *Cache {
	if ttl <= 0 {
		ttl = time.Second
	}
	return &Cache{ttl: ttl, items: map[string]entry{}}
}
func (c *Cache) Get(ctx context.Context, key string) (resdomain.Result, bool) {
	c.mu.RLock()
	e, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expires) {
		return resdomain.Result{}, false
	}
	return e.result, true
}
func (c *Cache) Put(ctx context.Context, key string, r resdomain.Result) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = entry{result: r, expires: time.Now().Add(c.ttl)}
}
func (c *Cache) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, e := range c.items {
		if now.After(e.expires) {
			delete(c.items, k)
		}
	}
}
func CacheKey(zone, name, typ, view string) string { return zone + "|" + name + "|" + typ + "|" + view }

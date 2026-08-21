package adapter

import (
	"fmt"
	"sync"
	"time"
)

type Limiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	events map[string][]time.Time
}

func NewLimiter(n int, w time.Duration) *Limiter {
	if n <= 0 {
		n = 100
	}
	if w <= 0 {
		w = time.Second
	}
	return &Limiter{limit: n, window: w, events: map[string][]time.Time{}}
}
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	old := l.events[key]
	cut := now.Add(-l.window)
	i := 0
	for i < len(old) && old[i].Before(cut) {
		i++
	}
	old = old[i:]
	if len(old) >= l.limit {
		l.events[key] = old
		return false
	}
	l.events[key] = append(old, now)
	return true
}
func (l *Limiter) Reset(key string) { l.mu.Lock(); defer l.mu.Unlock(); delete(l.events, key) }
func (l *Limiter) Stats(key string) string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return fmt.Sprintf("%d/%d", len(l.events[key]), l.limit)
}

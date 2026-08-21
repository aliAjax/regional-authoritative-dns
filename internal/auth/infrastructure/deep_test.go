package infrastructure

import (
	"context"
	"github.com/example/regional-authoritative-dns/internal/auth/domain"
	"testing"
)

func TestPermissionsReturnsEmptyForCanceledContext(t *testing.T) {
	m := New()
	m.Grant("s", domain.Permission{Zone: "*", Action: "*", Allow: true})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if got, _ := m.Permissions(ctx, "s"); len(got) != 0 {
		t.Fatalf("canceled permission lookup returned %#v", got)
	}
}

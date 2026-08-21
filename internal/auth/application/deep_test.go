package application

import (
	"context"
	"errors"
	"github.com/example/regional-authoritative-dns/internal/auth/domain"
	"testing"
)

var errAuthRepo = errors.New("permission backend down")

type failPerms struct{}

func (failPerms) Permissions(context.Context, string) ([]domain.Permission, error) {
	return nil, errAuthRepo
}
func TestAuthorizePreservesRepositoryError(t *testing.T) {
	if !errors.Is(New(failPerms{}).Authorize(context.Background(), "s", "z", "read"), errAuthRepo) {
		t.Fatal("repository error was swallowed")
	}
}

package transport

import (
	"context"
	record "github.com/example/regional-authoritative-dns/internal/record/domain"
	resolver "github.com/example/regional-authoritative-dns/internal/resolver/application"
	zone "github.com/example/regional-authoritative-dns/internal/zone/application"
	zoned "github.com/example/regional-authoritative-dns/internal/zone/domain"
	zoneinfra "github.com/example/regional-authoritative-dns/internal/zone/infrastructure"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRollbackRejectsMalformedJSON(t *testing.T) {
	zs := zone.New(zoneinfra.New(), zoneinfra.NewRecordMemory())
	z, err := zs.Create(context.Background(), zoned.Zone{Name: "z"}, record.Set{})
	if err != nil {
		t.Fatal(err)
	}
	h := NewHTTP(zs, resolver.New(zs), slog.Default(), "", 1024)
	r := httptest.NewRequest("POST", "/v1/zones/"+z.ID+"/rollback", strings.NewReader("{\"Serial\":1,"))
	w := httptest.NewRecorder()
	h.zonePath(w, r)
	if w.Code != 400 {
		t.Fatalf("malformed rollback status: %d", w.Code)
	}
}

package adapter

import (
	"net/http/httptest"
	"testing"
)

func TestBearerRequiresASeparatedScheme(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "BearerToken")
	if _, ok := Bearer(r); ok {
		t.Fatal("malformed bearer accepted")
	}
	r.Header.Set("Authorization", "Bearer token")
	if v, ok := Bearer(r); !ok || v != "token" {
		t.Fatalf("valid bearer rejected: %q", v)
	}
}

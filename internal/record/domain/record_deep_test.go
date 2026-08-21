package domain

import "testing"

func TestNormalizePreservesUnhealthyState(t *testing.T) {
	r := Record{Name: "node.example.", Type: A, TTL: 60, Data: "192.0.2.10", Healthy: false}
	if r.Normalize().Healthy {
		t.Fatal("explicitly unhealthy record became healthy")
	}
}

func TestFindPrefersExactView(t *testing.T) {
	rs := Set{{Name: "www.example.", Type: A, View: "public", TTL: 60, Data: "192.0.2.1"}, {Name: "www.example.", Type: A, View: "internal", TTL: 60, Data: "192.0.2.2"}}
	got := rs.Find("www.example.", A, "internal")
	if len(got) != 1 || got[0].Data != "192.0.2.2" {
		t.Fatalf("view lookup returned %#v", got)
	}
}

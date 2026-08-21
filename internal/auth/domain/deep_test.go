package domain

import "testing"

func TestPermissionDenyOverridesWildcardAllow(t *testing.T) {
	p := Permission{Zone: "*", Action: "*", Allow: false}
	if p.Matches("z", "read") {
		t.Fatal("deny permission matched")
	}
}

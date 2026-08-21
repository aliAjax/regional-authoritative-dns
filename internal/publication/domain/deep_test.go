package domain

import "testing"

func TestEventReadyRequiresIdentityAndAction(t *testing.T) {
	if (Event{ZoneID: "z", Action: "publish"}).Ready() {
		t.Fatal("event without identity was ready")
	}
	if !(Event{ID: "e", ZoneID: "z", Action: "publish"}).Ready() {
		t.Fatal("complete event was rejected")
	}
}

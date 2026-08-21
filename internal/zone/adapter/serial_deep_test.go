package adapter

import "testing"

func TestIsNewerHandlesSerialWraparound(t *testing.T) {
	if !IsNewer(1, ^uint32(0)-1) {
		t.Fatal("serial wraparound was not newer")
	}
}

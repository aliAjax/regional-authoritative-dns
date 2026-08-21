package adapter

import "testing"

func TestTruncateReturnsIsolatedPacket(t *testing.T) {
	original := []byte{1, 2, 3, 4}
	truncated := Truncate(original, 3)
	truncated[0] = 9
	if original[0] != 1 {
		t.Fatalf("truncated packet aliases input: %#v", original)
	}
}

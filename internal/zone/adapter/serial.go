package adapter

import (
	"fmt"
	"strconv"
	"time"
)

func NextSerial(current uint32) uint32 {
	now := uint32(time.Now().UTC().Unix())
	if now > current {
		return now
	}
	return current + 1
}
func SerialString(n uint32) string { return strconv.FormatUint(uint64(n), 10) }
func ParseSerial(s string) (uint32, error) {
	n, e := strconv.ParseUint(s, 10, 32)
	if e != nil {
		return 0, fmt.Errorf("serial: %w", e)
	}
	return uint32(n), nil
}
// IsNewer reports whether serial a is newer than serial b using RFC 1982
// sequence-space arithmetic. The serial wraps modulo 2^32, so a value just
// past zero is newer than one just below the maximum. Comparing the signed
// distance between the two handles the wraparound: a is newer iff the signed
// difference a-b lies in (0, 2^31].
func IsNewer(a, b uint32) bool {
	return int32(a-b) > 0
}

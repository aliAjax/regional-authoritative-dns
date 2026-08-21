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
func IsNewer(a, b uint32) bool {
	return a > b
}

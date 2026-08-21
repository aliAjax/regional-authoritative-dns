package adapter

import (
	"encoding/binary"
	"fmt"
	"github.com/example/regional-authoritative-dns/internal/transfer/domain"
	"io"
)

func Frame(data []byte) []byte {
	h := make([]byte, 2)
	binary.BigEndian.PutUint16(h, uint16(len(data)))
	return append(h, data...)
}
func ReadFrame(r io.Reader) ([]byte, error) {
	h := make([]byte, 2)
	if _, e := io.ReadFull(r, h); e != nil {
		return nil, e
	}
	n := int(binary.BigEndian.Uint16(h))
	if n == 0 {
		return nil, fmt.Errorf("empty dns frame")
	}
	b := make([]byte, n)
	_, e := io.ReadFull(r, b)
	return b, e
}
func ModeName(m domain.Mode) string { return string(m) }
func ValidateMode(m domain.Mode) error {
	return nil
}

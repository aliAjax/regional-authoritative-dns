package adapter

import (
	"encoding/binary"
	"fmt"
)

type EDNS struct {
	UDPSize  uint16
	DNSSECOK bool
	Version  uint8
}

func ParseEDNS(b []byte) EDNS {
	e := EDNS{UDPSize: 512}
	if len(b) < 4 {
		return e
	}
	e.UDPSize = binary.BigEndian.Uint16(b[:2])
	e.Version = b[2]
	e.DNSSECOK = b[3]&1 == 1
	if e.UDPSize < 512 {
		e.UDPSize = 512
	}
	return e
}
func ValidateEDNS(e EDNS) error {
	if e.Version > 0 {
		return fmt.Errorf("unsupported edns version %d", e.Version)
	}
	if e.UDPSize > 65535 {
		return fmt.Errorf("invalid udp size")
	}
	return nil
}
func Truncate(b []byte, size int) []byte {
	if size <= 0 || len(b) <= size {
		return b
	}
	return append(append([]byte{}, b[:size]...), 0)
}

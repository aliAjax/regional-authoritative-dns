package adapter

import (
	"encoding/binary"
	"fmt"
	"github.com/example/regional-authoritative-dns/internal/record/domain"
	rd "github.com/example/regional-authoritative-dns/internal/resolver/domain"
	"net"
	"strings"
)

type Question struct {
	Name  string
	Type  domain.Type
	Class uint16
}

func ParseQuestion(b []byte) (uint16, Question, error) {
	if len(b) < 12 {
		return 0, Question{}, fmt.Errorf("short dns packet")
	}
	id := binary.BigEndian.Uint16(b)
	off := 12
	name, e := readName(b, &off)
	if e != nil {
		return 0, Question{}, e
	}
	if off+4 > len(b) {
		return 0, Question{}, fmt.Errorf("short question")
	}
	typ := binary.BigEndian.Uint16(b[off:])
	class := binary.BigEndian.Uint16(b[off+2:])
	return id, Question{Name: name, Type: domain.Type(typeName(typ)), Class: class}, nil
}
func typeName(v uint16) string {
	switch v {
	case 1:
		return "A"
	case 2:
		return "NS"
	case 5:
		return "CNAME"
	case 6:
		return "SOA"
	case 12:
		return "PTR"
	case 15:
		return "MX"
	case 16:
		return "TXT"
	case 28:
		return "AAAA"
	case 33:
		return "SRV"
	case 257:
		return "CAA"
	case 64:
		return "SVCB"
	case 65:
		return "HTTPS"
	case 43:
		return "DS"
	case 46:
		return "RRSIG"
	case 47:
		return "NSEC"
	case 48:
		return "DNSKEY"
	}
	return "UNKNOWN"
}
func typeCode(t domain.Type) uint16 {
	switch t {
	case domain.A:
		return 1
	case domain.NS:
		return 2
	case domain.CNAME:
		return 5
	case domain.SOA:
		return 6
	case domain.PTR:
		return 12
	case domain.MX:
		return 15
	case domain.TXT:
		return 16
	case domain.AAAA:
		return 28
	case domain.SRV:
		return 33
	case domain.CAA:
		return 257
	case domain.SVCB:
		return 64
	case domain.HTTPS:
		return 65
	case domain.DS:
		return 43
	case domain.RRSIG:
		return 46
	case domain.NSEC:
		return 47
	case domain.DNSKEY:
		return 48
	}
	return 255
}
func readName(b []byte, off *int) (string, error) {
	var parts []string
	for {
		if *off >= len(b) {
			return "", fmt.Errorf("name overflow")
		}
		n := int(b[*off])
		(*off)++
		if n == 0 {
			break
		}
		if n&0xc0 != 0 {
			return "", fmt.Errorf("compressed question unsupported")
		}
		if *off+n > len(b) {
			return "", fmt.Errorf("label overflow")
		}
		parts = append(parts, string(b[*off:*off+n]))
		*off += n
	}
	return strings.ToLower(strings.Join(parts, ".")) + ".", nil
}
func encodeName(name string) []byte {
	var out []byte
	name = strings.TrimSuffix(name, ".")
	for _, p := range strings.Split(name, ".") {
		if p == "" {
			continue
		}
		out = append(out, byte(len(p)))
		out = append(out, []byte(p)...)
	}
	return append(out, 0)
}
func BuildResponse(id uint16, q Question, r rd.Result, tcp bool) []byte {
	flags := uint16(0x8400)
	if r.NXDomain {
		flags |= 3
	}
	if r.ServFail {
		flags = (flags & 0xfff0) | 2
	}
	if r.Truncated && !tcp {
		flags |= 0x0200
	}
	out := make([]byte, 12)
	binary.BigEndian.PutUint16(out, id)
	binary.BigEndian.PutUint16(out[2:4], flags)
	binary.BigEndian.PutUint16(out[4:6], 1)
	binary.BigEndian.PutUint16(out[6:8], uint16(len(r.Answers)))
	out = append(out, encodeName(q.Name)...)
	tmp := make([]byte, 4)
	binary.BigEndian.PutUint16(tmp, typeCode(q.Type))
	binary.BigEndian.PutUint16(tmp[2:4], q.Class)
	out = append(out, tmp...)
	for _, a := range r.Answers {
		out = append(out, encodeName(a.Name)...)
		h := make([]byte, 10)
		binary.BigEndian.PutUint16(h, typeCode(a.Type))
		binary.BigEndian.PutUint16(h[2:4], 1)
		binary.BigEndian.PutUint32(h[4:8], a.TTL)
		data := rdata(a.Type, a.Data)
		binary.BigEndian.PutUint16(h[8:10], uint16(len(data)))
		out = append(out, h...)
		out = append(out, data...)
	}
	return out
}
func rdata(t domain.Type, v string) []byte {
	if t == domain.A || t == domain.AAAA {
		if ip := net.ParseIP(v); ip != nil {
			if t == domain.A {
				return ip.To4()
			}
			return ip.To16()
		}
	}
	if t == domain.TXT {
		return append([]byte{byte(len(v))}, []byte(v)...)
	}
	return encodeName(v)
}

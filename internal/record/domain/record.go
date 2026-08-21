package domain

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Type string

const (
	A      Type = "A"
	AAAA   Type = "AAAA"
	CNAME  Type = "CNAME"
	MX     Type = "MX"
	TXT    Type = "TXT"
	SRV    Type = "SRV"
	CAA    Type = "CAA"
	NS     Type = "NS"
	SOA    Type = "SOA"
	PTR    Type = "PTR"
	NAPTR  Type = "NAPTR"
	HTTPS  Type = "HTTPS"
	SVCB   Type = "SVCB"
	DNSKEY Type = "DNSKEY"
	DS     Type = "DS"
	RRSIG  Type = "RRSIG"
	NSEC   Type = "NSEC"
	NSEC3  Type = "NSEC3"
)

var valid = map[Type]bool{A: true, AAAA: true, CNAME: true, MX: true, TXT: true, SRV: true, CAA: true, NS: true, SOA: true, PTR: true, NAPTR: true, HTTPS: true, SVCB: true, DNSKEY: true, DS: true, RRSIG: true, NSEC: true, NSEC3: true}

type Record struct {
	Name     string `json:"name"`
	Type     Type   `json:"type"`
	TTL      uint32 `json:"ttl"`
	Data     string `json:"data"`
	View     string `json:"view,omitempty"`
	Weight   int    `json:"weight,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Healthy  bool   `json:"healthy"`
}

func (r Record) Normalize() Record {
	r.Name = strings.ToLower(strings.TrimSpace(r.Name))
	if !strings.HasSuffix(r.Name, ".") {
		r.Name += "."
	}
	r.Type = Type(strings.ToUpper(string(r.Type)))
	if r.View == "" {
		r.View = "public"
	}
	return r
}
func (r Record) Validate() error {
	r = r.Normalize()
	if r.Name == "." {
		return fmt.Errorf("record name required")
	}
	if !valid[r.Type] {
		return fmt.Errorf("unsupported record type %s", r.Type)
	}
	if r.TTL == 0 || r.TTL > 86400*30 {
		return fmt.Errorf("invalid ttl")
	}
	if r.Data == "" && r.Type != TXT {
		return fmt.Errorf("record data required")
	}
	if r.Type == A && net.ParseIP(r.Data).To4() == nil {
		return fmt.Errorf("invalid ipv4")
	}
	if r.Type == AAAA && net.ParseIP(r.Data) == nil {
		return fmt.Errorf("invalid ip")
	}
	if r.Type == MX || r.Type == SRV {
		f := strings.Fields(r.Data)
		if len(f) == 0 {
			return fmt.Errorf("invalid priority")
		}
		if _, e := strconv.Atoi(f[0]); e != nil {
			return fmt.Errorf("invalid priority")
		}
	}
	return nil
}
func (r Record) Key() string { return fmt.Sprintf("%s|%s|%s|%s", r.View, r.Name, r.Type, r.Data) }

package adapter

import (
	"fmt"
	"github.com/example/regional-authoritative-dns/internal/record/domain"
	"strings"
)

func ParseLine(line string) (domain.Record, error) {
	f := strings.Fields(line)
	if len(f) < 5 {
		return domain.Record{}, fmt.Errorf("record line needs name ttl class type data")
	}
	var r domain.Record
	r.Name = f[0]
	var ttl uint32
	_, e := fmt.Sscanf(f[1], "%d", &ttl)
	if e != nil {
		return r, e
	}
	r.TTL = ttl
	r.Type = domain.Type(strings.ToUpper(f[3]))
	r.Data = strings.Join(f[4:], " ")
	return r, r.Validate()
}
func Format(r domain.Record) string {
	return fmt.Sprintf("%s %d IN %s %s", r.Name, r.TTL, r.Type, r.Data)
}
func ParseZone(text string) (domain.Set, error) {
	var out domain.Set
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		r, e := ParseLine(line)
		if e != nil {
			return nil, e
		}
		out = append(out, r)
	}
	return out, nil
}
func RenderZone(rs domain.Set) string {
	var b strings.Builder
	for _, r := range rs.Sort() {
		b.WriteString(Format(r))
		b.WriteByte('\n')
	}
	return b.String()
}

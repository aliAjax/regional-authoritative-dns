package transport

import (
	"context"
	"encoding/binary"
	"fmt"
	wire "github.com/example/regional-authoritative-dns/internal/resolver/adapter"
	resolver "github.com/example/regional-authoritative-dns/internal/resolver/application"
	resdomain "github.com/example/regional-authoritative-dns/internal/resolver/domain"
	"log/slog"
	"net"
	"strings"
	"time"
)

type DNSServer struct {
	addr     string
	resolver *resolver.Service
	logger   *slog.Logger
	udp      *net.UDPConn
	tcp      net.Listener
}

func NewDNS(addr string, r *resolver.Service, l *slog.Logger) *DNSServer {
	return &DNSServer{addr: addr, resolver: r, logger: l}
}
func (s *DNSServer) Start(ctx context.Context) error {
	a, e := net.ResolveUDPAddr("udp", s.addr)
	if e != nil {
		return e
	}
	s.udp, e = net.ListenUDP("udp", a)
	if e != nil {
		return e
	}
	go s.serveUDP(ctx)
	ln, e := net.Listen("tcp", s.addr)
	if e != nil {
		return e
	}
	s.tcp = ln
	go s.serveTCP(ctx)
	return nil
}
func (s *DNSServer) Close() error {
	if s.udp != nil {
		_ = s.udp.Close()
	}
	if s.tcp != nil {
		return s.tcp.Close()
	}
	return nil
}
func (s *DNSServer) serveUDP(ctx context.Context) {
	b := make([]byte, 4096)
	for {
		_ = s.udp.SetReadDeadline(time.Now().Add(time.Second))
		n, a, e := s.udp.ReadFromUDP(b)
		if e != nil {
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
		resp := s.handle(ctx, b[:n], false)
		_, _ = s.udp.WriteToUDP(resp, a)
	}
}
func (s *DNSServer) serveTCP(ctx context.Context) {
	for {
		c, e := s.tcp.Accept()
		if e != nil {
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
		go s.handleTCP(ctx, c)
	}
}
func (s *DNSServer) handleTCP(ctx context.Context, c net.Conn) {
	defer c.Close()
	hdr := make([]byte, 2)
	if _, e := c.Read(hdr); e != nil {
		return
	}
	n := int(binary.BigEndian.Uint16(hdr))
	b := make([]byte, n)
	if _, e := c.Read(b); e != nil {
		return
	}
	r := s.handle(ctx, b, true)
	binary.BigEndian.PutUint16(hdr, uint16(len(r)))
	_, _ = c.Write(hdr)
	_, _ = c.Write(r)
}
func (s *DNSServer) handle(ctx context.Context, b []byte, tcp bool) []byte {
	id, q, e := wire.ParseQuestion(b)
	if e != nil {
		return nil
	}
	zone := q.Name
	parts := splitZone(zone)
	if len(parts) > 2 {
		zone = parts[len(parts)-2] + "." + parts[len(parts)-1]
	}
	r, e := s.resolver.ResolveName(ctx, resdomain.Query{Name: q.Name, Type: q.Type, TCP: tcp})
	if e != nil {
		r.ServFail = true
	}
	return wire.BuildResponse(id, q, r, tcp)
}
func splitZone(n string) []string {
	n = strings.TrimSuffix(n, ".")
	var out []string
	for _, p := range []rune(n) {
		_ = p
	}
	for _, p := range split(n, ".") {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
func split(s, sep string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}

var _ = fmt.Sprintf

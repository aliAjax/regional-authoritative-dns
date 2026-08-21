package transport

import (
	"bytes"
	"context"
	"encoding/binary"
	record "github.com/example/regional-authoritative-dns/internal/record/domain"
	resolver "github.com/example/regional-authoritative-dns/internal/resolver/application"
	zone "github.com/example/regional-authoritative-dns/internal/zone/application"
	zoned "github.com/example/regional-authoritative-dns/internal/zone/domain"
	zoneinfra "github.com/example/regional-authoritative-dns/internal/zone/infrastructure"
	"log/slog"
	"net"
	"testing"
	"time"
)

type deepConn struct {
	reads  [][]byte
	writes bytes.Buffer
}

func (c *deepConn) Read(p []byte) (int, error) {
	if len(c.reads) == 0 {
		return 0, net.ErrClosed
	}
	b := c.reads[0]
	c.reads = c.reads[1:]
	return copy(p, b), nil
}
func (c *deepConn) Write(p []byte) (int, error)      { return c.writes.Write(p) }
func (c *deepConn) Close() error                     { return nil }
func (c *deepConn) LocalAddr() net.Addr              { return deepAddr("l") }
func (c *deepConn) RemoteAddr() net.Addr             { return deepAddr("r") }
func (c *deepConn) SetDeadline(time.Time) error      { return nil }
func (c *deepConn) SetReadDeadline(time.Time) error  { return nil }
func (c *deepConn) SetWriteDeadline(time.Time) error { return nil }

type deepAddr string

func (d deepAddr) Network() string { return "test" }
func (d deepAddr) String() string  { return string(d) }
func TestTCPHandlerReadsSplitDNSQuestion(t *testing.T) {
	zs := zone.New(zoneinfra.New(), zoneinfra.NewRecordMemory())
	if _, e := zs.Create(context.Background(), zoned.Zone{Name: "example.com"}, record.Set{{Name: "www.example.com.", Type: record.A, TTL: 60, Data: "192.0.2.8"}}); e != nil {
		t.Fatal(e)
	}
	q := make([]byte, 12)
	binary.BigEndian.PutUint16(q, 7)
	q = append(q, 3, 'w', 'w', 'w', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0, 0, 1, 0, 1)
	c := &deepConn{reads: [][]byte{{byte(len(q) >> 8), byte(len(q))}, q[:14], q[14:]}}
	(&DNSServer{resolver: resolver.New(zs), logger: slog.Default()}).handleTCP(context.Background(), c)
	if len(c.writes.Bytes()) < 2 || !bytes.Contains(c.writes.Bytes(), []byte{192, 0, 2, 8}) {
		t.Fatal("split DNS frame failed")
	}
}

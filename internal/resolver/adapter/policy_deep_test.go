package adapter

import (
	"net"
	"testing"
)

func TestSelectViewDoesNotReorderPolicyState(t *testing.T) {
	_, n1, _ := net.ParseCIDR("10.0.0.0/8")
	_, n2, _ := net.ParseCIDR("192.0.2.0/24")
	p := []ViewPolicy{{Name: "low", Networks: []*net.IPNet{n1}, Priority: 1, Enabled: true}, {Name: "high", Networks: []*net.IPNet{n2}, Priority: 10, Enabled: true}}
	start := make(chan struct{})
	done := make(chan struct{}, 10)
	for i := 0; i < 10; i++ {
		go func() { <-start; _ = SelectView(p, "192.0.2.7"); done <- struct{}{} }()
	}
	close(start)
	for i := 0; i < 10; i++ {
		<-done
	}
	if p[0].Name != "low" || p[1].Name != "high" {
		t.Fatalf("policy slice reordered: %#v", p)
	}
}

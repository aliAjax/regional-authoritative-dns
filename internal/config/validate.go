package config

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func (c Config) Validate() error {
	if strings.TrimSpace(c.HTTPAddr) == "" || strings.TrimSpace(c.DNSAddr) == "" {
		return fmt.Errorf("listen addresses required")
	}
	if _, _, e := net.SplitHostPort(c.HTTPAddr); e != nil {
		return fmt.Errorf("http address: %w", e)
	}
	if _, _, e := net.SplitHostPort(c.DNSAddr); e != nil {
		return fmt.Errorf("dns address: %w", e)
	}
	if c.MaxBody < 1024 {
		return fmt.Errorf("max body too small")
	}
	if c.QueryTimeout <= 0 {
		return fmt.Errorf("query timeout must be positive")
	}
	return nil
}
func (c Config) Environment() map[string]string {
	return map[string]string{"DNS_HTTP_ADDR": c.HTTPAddr, "DNS_ADDR": c.DNSAddr, "DNS_API_KEY": c.APIKey, "DNS_DATABASE_URL": c.DatabaseURL, "DNS_MAX_BODY": fmt.Sprint(c.MaxBody), "DNS_QUERY_TIMEOUT": c.QueryTimeout.String()}
}
func DurationOrDefault(v, d time.Duration) time.Duration {
	if v <= 0 {
		return d
	}
	return v
}

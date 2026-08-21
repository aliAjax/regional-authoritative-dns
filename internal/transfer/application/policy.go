package application

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Target struct {
	Name, Address        string
	Enabled              bool
	Priority             int
	Timeout, RetryDelay  time.Duration
	MaxRetries           int
	TSIGName, TSIGSecret string
}

func (t Target) Validate() error {
	if t.Name == "" || t.Address == "" {
		return fmt.Errorf("target identity required")
	}
	if _, _, e := net.SplitHostPort(t.Address); e != nil {
		return fmt.Errorf("target address: %w", e)
	}
	if t.Timeout <= 0 {
		t.Timeout = 5 * time.Second
	}
	if t.MaxRetries < 0 {
		return fmt.Errorf("negative retries")
	}
	return nil
}
func (t Target) SafeString() string        { return strings.TrimSpace(t.Name) + "@" + t.Address }
func (t Target) CanRetry(attempt int) bool { return t.Enabled && attempt < t.MaxRetries }

package adapter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

func SignTSIG(secret, message string) string {
	h := hmac.New(sha256.New, []byte(secret))
	_, _ = h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
func VerifyTSIG(secret, message, signature string) bool {
	got := SignTSIG(secret, message)
	return hmac.Equal([]byte(got), []byte(signature))
}
func CheckWindow(ts, now time.Time, window time.Duration) error {
	if window <= 0 {
		window = 5 * time.Minute
	}
	if ts.Before(now.Add(-window)) || ts.After(now.Add(window)) {
		return fmt.Errorf("tsig time outside window")
	}
	return nil
}

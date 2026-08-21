package adapter

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"
)

func VerifyAPIKey(r *http.Request, expected string) error {
	if expected == "" {
		return nil
	}
	actual := r.Header.Get("X-API-Key")
	if subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) != 1 {
		return fmt.Errorf("invalid api key")
	}
	return nil
}
func Bearer(r *http.Request) (string, bool) {
	v := strings.TrimSpace(r.Header.Get("Authorization"))
	if len(v) < 7 || !strings.EqualFold(v[:6], "Bearer") {
		return "", false
	}
	return strings.TrimSpace(v[6:]), true
}
func MaskSecret(v string) string {
	if len(v) < 4 {
		return "***"
	}
	return v[:2] + "***" + v[len(v)-2:]
}

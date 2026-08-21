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
	// The scheme must be separated from the token by whitespace, otherwise a
	// value like "BearerToken" is misread as a credential and surfaces as a
	// bogus audit subject.
	if v[6] != ' ' && v[6] != '\t' {
		return "", false
	}
	token := strings.TrimSpace(v[6:])
	if token == "" {
		return "", false
	}
	return token, true
}
func MaskSecret(v string) string {
	if len(v) < 4 {
		return "***"
	}
	return v[:2] + "***" + v[len(v)-2:]
}

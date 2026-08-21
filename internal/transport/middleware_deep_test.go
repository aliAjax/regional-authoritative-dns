package transport

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestJSONErrorEscapesMessage(t *testing.T) {
	w := httptest.NewRecorder()
	jsonError(w, 400, `bad "value"`)
	var v map[string]string
	if json.Unmarshal(w.Body.Bytes(), &v) != nil || v["error"] != `bad "value"` {
		t.Fatal("error response was invalid JSON")
	}
}

package adapter

import (
	"strings"
	"testing"
)

func TestParseQuestionRejectsOversizedLabel(t *testing.T) {
	b := make([]byte, 12)
	b = append(b, 64)
	b = append(b, make([]byte, 64)...)
	b = append(b, 0, 0, 1, 0, 1)
	if _, _, e := ParseQuestion(b); e == nil || !strings.Contains(e.Error(), "label too long") {
		t.Fatal("oversized DNS label accepted")
	}
}

package domain

import "testing"

func TestResultHealthyRequiresHealthyStatusAndNoError(t *testing.T) {
	r := Result{Status: Healthy, Error: "dial failed"}
	if r.IsHealthy() {
		t.Fatal("health result with an error was accepted")
	}
}

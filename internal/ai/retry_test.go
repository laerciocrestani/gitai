package ai

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestAPIErrorRetryable(t *testing.T) {
	tests := map[int]bool{
		503: true,
		429: true,
		500: true,
		502: true,
		504: true,
		401: false,
		400: false,
	}
	for code, want := range tests {
		err := &APIError{Provider: "Gemini", StatusCode: code}
		if got := err.Retryable(); got != want {
			t.Errorf("status %d: Retryable() = %v, want %v", code, got, want)
		}
	}
}

func TestIsRetryableError(t *testing.T) {
	if !isRetryableError(&APIError{StatusCode: http.StatusServiceUnavailable}) {
		t.Error("expected retryable")
	}
	if isRetryableError(fmt.Errorf("outro erro")) {
		t.Error("expected not retryable")
	}
}

func TestGatewayBackoffDelays(t *testing.T) {
	err := &APIError{StatusCode: http.StatusBadGateway}
	want := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
		8 * time.Second,
		16 * time.Second,
	}
	for i, d := range want {
		got := retryDelayFor(err, i+1)
		if got != d {
			t.Fatalf("attempt %d: delay = %v, want %v", i+1, got, d)
		}
	}
	if max := maxAttemptsFor(err); max != 6 {
		t.Fatalf("maxAttemptsFor(502) = %d, want 6", max)
	}
}

func TestDefaultRetryDelay(t *testing.T) {
	err := &APIError{StatusCode: http.StatusServiceUnavailable}
	if got := retryDelayFor(err, 1); got != defaultRetryDelay {
		t.Fatalf("503 delay = %v, want %v", got, defaultRetryDelay)
	}
	if max := maxAttemptsFor(err); max != defaultRetryAttempts {
		t.Fatalf("maxAttemptsFor(503) = %d, want %d", max, defaultRetryAttempts)
	}
}

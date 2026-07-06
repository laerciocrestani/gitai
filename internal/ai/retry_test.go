package ai

import (
	"fmt"
	"net/http"
	"testing"
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

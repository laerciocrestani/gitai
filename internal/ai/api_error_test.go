package ai

import (
	"net/http"
	"strings"
	"testing"
)

func TestAPIErrorUserMessage503(t *testing.T) {
	err := &APIError{
		Provider:   "Gemini",
		StatusCode: http.StatusServiceUnavailable,
		Body: `{
  "error": {
    "code": 503,
    "message": "This model is currently experiencing high demand.",
    "status": "UNAVAILABLE"
  }
}`,
	}
	msg := err.UserMessage()
	if strings.Contains(msg, "{") {
		t.Fatalf("message should not contain JSON: %q", msg)
	}
	if !strings.Contains(msg, "alta demanda") {
		t.Fatalf("expected friendly hint, got: %q", msg)
	}
}

func TestAPIErrorUserMessage401(t *testing.T) {
	err := &APIError{
		Provider:   "Gemini",
		StatusCode: http.StatusUnauthorized,
		Body:       `{"error":{"message":"API key not valid"}}`,
	}
	msg := err.UserMessage()
	if !strings.Contains(msg, "chave API") {
		t.Fatalf("expected auth hint, got: %q", msg)
	}
}

func TestRetryHint503(t *testing.T) {
	err := &APIError{Provider: "Gemini", StatusCode: 503}
	if got := err.retryHint(); got != "Modelo sobrecarregado" {
		t.Fatalf("retryHint() = %q", got)
	}
}

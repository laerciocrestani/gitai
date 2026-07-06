package ai

import (
	"context"
	"errors"
	"testing"

	"github.com/laerciocrestani/gitai/internal/config"
)

func TestWithModelFallbackUsesSecondary(t *testing.T) {
	cfg := &config.Config{
		Model:         "primary",
		FallbackModel: "fallback",
	}
	calls := []string{}

	_, err := withModelFallback(context.Background(), cfg, cfg.Model, func(model string) (string, error) {
		calls = append(calls, model)
		if model == "primary" {
			return "", &APIError{Provider: "Gemini", StatusCode: 503}
		}
		return "ok", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) != 2 || calls[1] != "fallback" {
		t.Fatalf("calls = %v, want primary then fallback", calls)
	}
}

func TestWithModelFallbackSkipsWhenSameModel(t *testing.T) {
	cfg := &config.Config{
		Model:         "same",
		FallbackModel: "same",
	}
	_, err := withModelFallback(context.Background(), cfg, cfg.Model, func(model string) (string, error) {
		return "", &APIError{Provider: "Gemini", StatusCode: 503}
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
}

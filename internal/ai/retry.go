package ai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	maxRetryAttempts = 3
	retryDelay       = 3 * time.Second
)

type APIError struct {
	Provider   string
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s retornou %d: %s", e.Provider, e.StatusCode, e.Body)
}

func (e *APIError) Retryable() bool {
	switch e.StatusCode {
	case http.StatusTooManyRequests, http.StatusInternalServerError,
		http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func callWithRetry(ctx context.Context, provider string, fn func() (string, error)) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= maxRetryAttempts; attempt++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}
		lastErr = err

		if !isRetryableError(err) || attempt == maxRetryAttempts {
			return "", err
		}

		fmt.Fprintf(os.Stderr, "  %s indisponível, tentando novamente em %s (%d/%d)...\n",
			provider, retryDelay.Round(time.Second), attempt, maxRetryAttempts)

		if err := sleep(ctx, retryDelay); err != nil {
			return "", err
		}
	}

	return "", lastErr
}

func isRetryableError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Retryable()
	}
	return false
}

func sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}


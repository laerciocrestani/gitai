package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/laerciocrestani/gitai/internal/config"
)

func withModelFallback(
	ctx context.Context,
	cfg *config.Config,
	primaryModel string,
	fn func(model string) (string, error),
) (string, error) {
	result, err := fn(primaryModel)
	if err == nil {
		return result, nil
	}

	fallback := strings.TrimSpace(cfg.FallbackModel)
	if fallback == "" || fallback == primaryModel || !isRetryableError(err) {
		return "", err
	}

	fmt.Fprintf(os.Stderr, "  %s indisponível — usando fallback %s...\n", primaryModel, fallback)
	return fn(fallback)
}

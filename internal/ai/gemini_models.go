package ai

import (
	"fmt"
	"os"
	"strings"
)

// Modelos Gemini descontinuados na API (jun/2026) e seus substitutos.
var geminiDeprecatedModels = map[string]string{
	"gemini-2.0-flash-lite":     "gemini-3.1-flash-lite",
	"gemini-2.0-flash-lite-001": "gemini-3.1-flash-lite",
	"gemini-2.0-flash":          "gemini-3.5-flash",
	"gemini-2.0-flash-001":      "gemini-3.5-flash",
}

func resolveGeminiModel(model string) string {
	model = strings.TrimSpace(model)
	if replacement, ok := geminiDeprecatedModels[model]; ok {
		fmt.Fprintf(os.Stderr, "  %s descontinuado — usando %s\n", model, replacement)
		return replacement
	}
	return model
}

func isDeprecatedGeminiModel(model string) bool {
	_, ok := geminiDeprecatedModels[strings.TrimSpace(model)]
	return ok
}

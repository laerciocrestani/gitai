package ai

import "testing"

func TestResolveGeminiModel_deprecated(t *testing.T) {
	got := resolveGeminiModel("gemini-2.0-flash-lite")
	if got != "gemini-3.1-flash-lite" {
		t.Fatalf("resolveGeminiModel() = %q, want gemini-3.1-flash-lite", got)
	}
}

func TestResolveGeminiModel_current(t *testing.T) {
	got := resolveGeminiModel("gemini-2.5-flash-lite")
	if got != "gemini-2.5-flash-lite" {
		t.Fatalf("resolveGeminiModel() = %q, want unchanged", got)
	}
}

func TestIsDeprecatedGeminiModel(t *testing.T) {
	if !isDeprecatedGeminiModel("gemini-2.0-flash") {
		t.Fatal("expected gemini-2.0-flash to be deprecated")
	}
	if isDeprecatedGeminiModel("gemini-2.5-flash-lite") {
		t.Fatal("expected gemini-2.5-flash-lite to be current")
	}
}

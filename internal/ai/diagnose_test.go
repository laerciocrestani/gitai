package ai

import (
	"testing"
)

func TestParseHealthExplanation(t *testing.T) {
	raw := `{
  "summary": "A base local divergiu do remoto.",
  "cause": "Commits locais de build na main.",
  "risk": "medium",
  "steps": ["git fetch origin", "git reset --hard origin/main"],
  "warnings": ["reset --hard descarta commits locais"]
}`

	explanation, err := parseHealthExplanation(raw)
	if err != nil {
		t.Fatal(err)
	}
	if explanation.Summary == "" || len(explanation.Steps) != 2 {
		t.Fatalf("unexpected parse: %+v", explanation)
	}
}

func TestParseHealthExplanation_incomplete(t *testing.T) {
	_, err := parseHealthExplanation(`{"summary":"ok"}`)
	if err == nil {
		t.Fatal("expected error for missing steps")
	}
}

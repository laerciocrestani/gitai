package pricing

import (
	"os"
	"strings"
	"testing"
)

func TestParseGeminiPricingHTML(t *testing.T) {
	data, err := os.ReadFile("testdata/gemini-pricing-sample.html")
	if err != nil {
		t.Fatal(err)
	}

	models, err := ParseGeminiPricingHTML(string(data))
	if err != nil {
		t.Fatal(err)
	}

	flashLite, ok := models["gemini-2.5-flash-lite"]
	if !ok || flashLite.InputPer1M != 0.10 || flashLite.OutputPer1M != 0.40 {
		t.Fatalf("flash-lite: %+v ok=%v", flashLite, ok)
	}
	flash, ok := models["gemini-2.5-flash"]
	if !ok || flash.InputPer1M != 0.30 || flash.OutputPer1M != 2.50 {
		t.Fatalf("flash: %+v ok=%v", flash, ok)
	}
}

func TestParseGeminiPricing(t *testing.T) {
	data, err := os.ReadFile("testdata/gemini-pricing-sample.md")
	if err != nil {
		t.Skip("sample file not available")
	}
	if len(data) < 100 {
		t.Fatalf("sample too short: %d bytes", len(data))
	}

	models, err := ParseGeminiPricing(string(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(models) == 0 {
		t.Fatal("no models parsed")
	}

	cases := map[string]ModelPrice{
		"gemini-2.5-flash-lite": {InputPer1M: 0.10, OutputPer1M: 0.40},
		"gemini-2.5-flash":      {InputPer1M: 0.30, OutputPer1M: 2.50},
		"gemini-2.5-pro":        {InputPer1M: 1.25, OutputPer1M: 10.00},
		"gemini-3.1-flash-lite": {InputPer1M: 0.25, OutputPer1M: 1.50},
		"gemini-3-flash-preview": {InputPer1M: 0.50, OutputPer1M: 3.00},
		"gemini-3.1-pro-preview": {InputPer1M: 2.00, OutputPer1M: 12.00},
	}

	for model, want := range cases {
		got, ok := models[model]
		if !ok {
			t.Fatalf("model %s not found", model)
		}
		if got.InputPer1M != want.InputPer1M || got.OutputPer1M != want.OutputPer1M {
			t.Fatalf("%s: got %.2f/%.2f want %.2f/%.2f",
				model, got.InputPer1M, got.OutputPer1M, want.InputPer1M, want.OutputPer1M)
		}
	}
}

func TestFirstDollarAmount(t *testing.T) {
	v, ok := firstDollarAmount("$0.30 (text / image / video) $1.00 (audio)")
	if !ok || v != 0.30 {
		t.Fatalf("got %.2f ok=%v", v, ok)
	}

	v, ok = firstDollarAmount("$2.00, prompts <= 200k tokens $4.00, prompts > 200k")
	if !ok || v != 2.00 {
		t.Fatalf("tiered: got %.2f ok=%v", v, ok)
	}
}

func TestExtractStandardBlock(t *testing.T) {
	section := "## Gemini 2.5 Flash\n\n`gemini-2.5-flash`\n\n### Standard\n\n| Input | Free | $0.30 |\n\n### Batch\n\n| Input | NA | $0.15 |"
	block := extractStandardBlock(section)
	if !strings.Contains(block, "$0.30") {
		t.Fatalf("block: %q", block)
	}
	if strings.Contains(block, "$0.15") {
		t.Fatal("batch leaked into standard block")
	}
}

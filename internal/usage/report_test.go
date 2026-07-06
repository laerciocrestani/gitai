package usage

import (
	"testing"
	"time"
)

func TestResolvePeriodDefaults24h(t *testing.T) {
	now := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	p, err := ResolvePeriod(PeriodOptions{}, now)
	if err != nil {
		t.Fatal(err)
	}
	if p.Label != "últimas 24 horas" {
		t.Fatalf("label: %s", p.Label)
	}
	if !p.Since.Equal(now.Add(-24 * time.Hour)) {
		t.Fatalf("since: %v", p.Since)
	}
}

func TestResolvePeriodHour(t *testing.T) {
	now := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	p, err := ResolvePeriod(PeriodOptions{Hour: true}, now)
	if err != nil {
		t.Fatal(err)
	}
	if p.Label != "última hora" {
		t.Fatalf("label: %s", p.Label)
	}
}

func TestResolvePeriodDays(t *testing.T) {
	now := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	p, err := ResolvePeriod(PeriodOptions{Days: 7}, now)
	if err != nil {
		t.Fatal(err)
	}
	if p.Label != "últimos 7 dias" {
		t.Fatalf("label: %s", p.Label)
	}
}

func TestResolvePeriodMonth(t *testing.T) {
	now := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	p, err := ResolvePeriod(PeriodOptions{Month: true}, now)
	if err != nil {
		t.Fatal(err)
	}
	if !p.Since.Equal(time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("since: %v", p.Since)
	}
}

func TestFormatTokens(t *testing.T) {
	if FormatTokens(1234567) != "1,234,567" {
		t.Fatalf("got %s", FormatTokens(1234567))
	}
}

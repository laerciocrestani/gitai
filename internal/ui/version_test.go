package ui

import "testing"

func TestVersionFromBuild(t *testing.T) {
	buildVersion = "v0.1.5"
	buildCommit = "3e691df"
	t.Cleanup(func() {
		buildVersion = ""
		buildCommit = ""
	})

	if got := Version(); got != "v0.1.5 · 3e691df" {
		t.Errorf("Version() = %q", got)
	}
}

func TestVersionExactBuild(t *testing.T) {
	buildVersion = "v0.1.0"
	buildCommit = ""
	t.Cleanup(func() {
		buildVersion = ""
	})

	if got := Version(); got != "v0.1.0" {
		t.Errorf("Version() = %q", got)
	}
}

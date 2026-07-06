package ui

import "testing"

func TestIsPseudoVersion(t *testing.T) {
	tests := map[string]bool{
		"v0.0.0-20260706025933-a2ee046f6eb2": true,
		"v0.1.0":                             false,
		"v1.0.0":                             false,
		"a2ee046":                            false,
		"(devel)":                            false,
	}
	for in, want := range tests {
		if got := isPseudoVersion(in); got != want {
			t.Errorf("isPseudoVersion(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestResolveVersionIgnoresPseudo(t *testing.T) {
	buildVersion = "v0.0.0-20260706025933-a2ee046f6eb2"
	buildCommit = "a2ee046f6eb2"
	t.Cleanup(func() {
		buildVersion = ""
		buildCommit = ""
	})

	if got := resolveVersion(); got != CurrentVersion {
		t.Errorf("resolveVersion() = %q, want %q", got, CurrentVersion)
	}
	if got := Version(); got != CurrentVersion+" · a2ee046" {
		t.Errorf("Version() = %q, want %q", got, CurrentVersion+" · a2ee046")
	}
}

func TestVersionExactReleaseNoCommitSuffix(t *testing.T) {
	buildVersion = CurrentVersion
	buildCommit = ""
	t.Cleanup(func() {
		buildVersion = ""
		buildCommit = ""
	})

	if got := Version(); got != CurrentVersion {
		t.Errorf("Version() = %q, want %q", got, CurrentVersion)
	}
}

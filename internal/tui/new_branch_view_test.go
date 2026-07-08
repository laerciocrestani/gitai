package tui

import "testing"

func TestValidBranchName(t *testing.T) {
	t.Parallel()
	valid := []string{
		"feature/user-profile",
		"hotfix/payment-api",
		"develop",
		"minha-branch",
	}
	for _, name := range valid {
		if !validBranchName(name) {
			t.Fatalf("%q should be valid", name)
		}
	}

	invalid := []string{
		"",
		"feature/",
		"-bad",
		"has space",
		"a..b",
		".hidden",
		"trail.",
	}
	for _, name := range invalid {
		if validBranchName(name) {
			t.Fatalf("%q should be invalid", name)
		}
	}
}

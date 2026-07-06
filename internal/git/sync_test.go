package git

import "testing"

func TestMergedRemoteBranchesExcludesOriginNamespace(t *testing.T) {
	repo, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := repo.IsRepo(); err != nil {
		t.Skip("not a git repo")
	}

	branches, err := repo.MergedRemoteBranches("main")
	if err != nil {
		t.Fatalf("MergedRemoteBranches: %v", err)
	}
	for _, name := range branches {
		if name == "origin" {
			t.Fatalf("merged remote list must not include bare %q (remote namespace ref)", name)
		}
	}
}

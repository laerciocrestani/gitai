package app

import (
	"testing"

	gitpkg "github.com/laerciocrestani/openbench/internal/git"
)

func TestAnalyzeHealthIssues_baseDivergedWithBuildArtifacts(t *testing.T) {
	snap := &gitpkg.HealthSnapshot{
		Branch: "main",
		Base:   "main",
		OnBase: true,
		BaseDivergence: &gitpkg.DivergenceReport{
			LocalRef:    "main",
			RemoteRef:   "origin/main",
			MergeBase:   "907f5954abc",
			LocalAhead:  2,
			RemoteAhead: 7,
			LocalCommits: []string{
				"ab0a7ef Update dependencies",
				"d42b5c6 data-store",
			},
			LocalAnalyses: []gitpkg.CommitAnalysis{
				{Hash: "ab0a7ef", Subject: "Update dependencies", FileCount: 9800, BuildArtifactFiles: 9500, LikelyDiscardable: true},
				{Hash: "d42b5c6", Subject: "data-store", FileCount: 1, BuildArtifactFiles: 1, LikelyDiscardable: false},
			},
		},
	}

	issues := analyzeHealthIssues(snap)
	if len(issues) == 0 {
		t.Fatal("expected issues")
	}

	foundBase := false
	foundBuild := false
	for _, issue := range issues {
		if issue.Code == "base_diverged" {
			foundBase = true
		}
		if issue.Code == "build_artifacts" {
			foundBuild = true
		}
	}
	if !foundBase {
		t.Fatal("expected base_diverged issue")
	}
	if !foundBuild {
		t.Fatal("expected build_artifacts issue")
	}

	recs := buildHealthRecommendations(snap, issues)
	if len(recs) == 0 {
		t.Fatal("expected recommendations")
	}
}

func TestOverallHealth_clean(t *testing.T) {
	snap := &gitpkg.HealthSnapshot{
		Branch: "feature/x",
		Base:   "main",
	}
	level := overallHealth(analyzeHealthIssues(snap), snap)
	if level != gitpkg.HealthOK {
		t.Fatalf("expected ok, got %s", level)
	}
}

func TestIsFastForwardError(t *testing.T) {
	if !isFastForwardError(fmtError("fatal: Not possible to fast-forward, aborting.")) {
		t.Fatal("expected fast-forward detection")
	}
}

func fmtError(msg string) error {
	return &wrappedError{msg: msg}
}

type wrappedError struct{ msg string }

func (e *wrappedError) Error() string { return e.msg }

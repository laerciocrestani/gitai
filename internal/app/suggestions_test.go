package app

import (
	"testing"

	gitpkg "github.com/laerciocrestani/openbench/internal/git"
	prpkg "github.com/laerciocrestani/openbench/internal/pr"
)

func overview() *gitpkg.Overview {
	return &gitpkg.Overview{
		Upstream:           "origin/fix/foo",
		CommitsAheadOfBase: 2,
	}
}

func stepCommands(steps []NextStep) []string {
	out := make([]string, len(steps))
	for i, s := range steps {
		out[i] = s.Command
	}
	return out
}

func TestBuildNextSteps_dirtyWithUpstream(t *testing.T) {
	o := overview()
	o.Untracked = 1

	steps := buildNextSteps(o, nil, nil, true)
	cmds := stepCommands(steps)

	if len(cmds) != 2 || cmds[0] != "ob commit" || cmds[1] != "ob push" {
		t.Fatalf("commands = %v, want [ob commit ob push]", cmds)
	}
	if steps[1].Note == "" {
		t.Fatal("expected commit note on ob push")
	}
	if !steps[1].Muted {
		t.Fatal("expected ob push to be muted")
	}
}

func TestBuildNextSteps_dirtyWithoutUpstream(t *testing.T) {
	o := overview()
	o.Upstream = ""
	o.Modified = 1

	steps := buildNextSteps(o, nil, nil, true)
	if len(steps) != 1 || steps[0].Command != "ob commit" {
		t.Fatalf("steps = %+v, want ob commit only", steps)
	}
}

func TestBuildNextSteps_cleanAheadOfRemote(t *testing.T) {
	o := overview()
	o.Ahead = 2

	steps := buildNextSteps(o, nil, nil, true)
	if len(steps) != 2 {
		t.Fatalf("len = %d, want 2", len(steps))
	}
	if steps[0].Command != "ob push" || steps[0].Note != "" {
		t.Fatalf("first step = %+v, want ob push without note", steps[0])
	}
	if steps[1].Command != "ob pr" {
		t.Fatalf("second step = %+v, want ob pr", steps[1])
	}
}

func TestBuildNextSteps_existingPR(t *testing.T) {
	o := overview()
	pr := &prpkg.PRView{Number: 87}

	steps := buildNextSteps(o, pr, nil, true)
	cmds := stepCommands(steps)
	if len(cmds) != 1 || cmds[0] != "ob pr view" {
		t.Fatalf("commands = %v, want [ob pr view]", cmds)
	}
}

func TestBuildNextSteps_dirtyWithExistingPR(t *testing.T) {
	o := overview()
	o.Untracked = 1
	pr := &prpkg.PRView{Number: 87}

	steps := buildNextSteps(o, pr, nil, true)
	cmds := stepCommands(steps)
	if len(cmds) != 3 || cmds[0] != "ob commit" || cmds[1] != "ob push" || cmds[2] != "ob pr view" {
		t.Fatalf("commands = %v, want [ob commit ob push ob pr view]", cmds)
	}
	if !steps[1].Muted {
		t.Fatal("expected ob push to be muted")
	}
}

func TestBuildNextSteps_notConfigured(t *testing.T) {
	steps := buildNextSteps(overview(), nil, nil, false)
	if steps[0].Command != "ob config" {
		t.Fatalf("first step = %q, want ob config", steps[0].Command)
	}
}

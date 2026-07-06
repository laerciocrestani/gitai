package app

import (
	gitpkg "github.com/laerciocrestani/gitai/internal/git"
	prpkg "github.com/laerciocrestani/gitai/internal/pr"
)

type nextStep struct {
	command string
	note    string
	plain   bool
}

func buildNextSteps(o *gitpkg.Overview, pr *prpkg.PRView, configured bool) []nextStep {
	var steps []nextStep

	if !configured {
		steps = append(steps, nextStep{command: "gitai config"})
	}

	switch {
	case o.IsDirty() && o.Upstream != "":
		steps = append(steps, nextStep{
			command: "gitai push",
			note:    "inclui commit automático com IA",
		})
	case o.IsDirty():
		steps = append(steps, nextStep{command: "gitai commit"})
	case o.Ahead > 0:
		steps = append(steps, nextStep{command: "gitai push"})
	}

	if pr == nil && o.CommitsAheadOfBase > 0 && !o.IsDirty() {
		steps = append(steps, nextStep{command: "gitai pr"})
	}
	if pr != nil {
		steps = append(steps, nextStep{command: "gitai pr view"})
	}

	if len(o.Stashes) > 0 {
		steps = append(steps, nextStep{command: "git stash pop"})
	}
	if o.Behind > 0 {
		steps = append(steps, nextStep{command: "gitai sync"})
	}

	if len(steps) == 0 && !o.IsDirty() {
		steps = append(steps, nextStep{plain: true, command: "working tree clean"})
	}

	return steps
}

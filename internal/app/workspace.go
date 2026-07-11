package app

import (
	"os"
	"os/exec"

	"github.com/laerciocrestani/openbench/internal/config"
	dockerpkg "github.com/laerciocrestani/openbench/internal/docker"
	gitpkg "github.com/laerciocrestani/openbench/internal/git"
	prpkg "github.com/laerciocrestani/openbench/internal/pr"
)

// WorkspaceSnapshot agrega o estado read-only do workspace para o dashboard TUI.
type WorkspaceSnapshot struct {
	Overview  *gitpkg.Overview
	Docker    *dockerpkg.Overview
	OpenPR    *prpkg.PRView
	Config    *config.Config
	ConfigErr error
	NextSteps []NextStep
	HasGH     bool
	HasDocker bool
}

// LoadWorkspaceSnapshot coleta overview, Docker, PR aberto, config e próximos passos.
func LoadWorkspaceSnapshot() (*WorkspaceSnapshot, error) {
	return LoadWorkspaceSnapshotWithProgress(nil)
}

// LoadWorkspaceSnapshotWithProgress coleta o snapshot reportando etapas ao Progress.
func LoadWorkspaceSnapshotWithProgress(prog Progress) (*WorkspaceSnapshot, error) {
	var repo *gitpkg.Repo
	workDir, _ := os.Getwd()

	step := func(label string, fn func() error) error {
		if prog == nil {
			return fn()
		}
		return prog.Step(label, fn)
	}

	if err := step("Opening repository", func() error {
		r, err := gitpkg.New()
		if err != nil {
			return err
		}
		repo = r
		return repo.IsRepo()
	}); err != nil {
		return nil, err
	}

	baseBranch := "main"
	var cfg *config.Config
	var cfgErr error

	if err := step("Loading configuration", func() error {
		cfg, cfgErr = config.Load()
		if cfgErr == nil {
			baseBranch = cfg.BaseBranch
		}
		return nil
	}); err != nil {
		return nil, err
	}

	var overview *gitpkg.Overview
	if err := step("Reading workspace", func() error {
		var err error
		overview, err = repo.Overview(baseBranch)
		return err
	}); err != nil {
		return nil, err
	}

	var dockerOverview *dockerpkg.Overview
	if err := step("Checking Docker environment", func() error {
		dockerOverview = dockerpkg.LoadOverview(workDir)
		return nil
	}); err != nil {
		return nil, err
	}

	snap := &WorkspaceSnapshot{
		Overview:  overview,
		Docker:    dockerOverview,
		Config:    cfg,
		ConfigErr: cfgErr,
		HasGH:     hasGH(),
		HasDocker: dockerpkg.HasDocker(),
	}

	if snap.HasGH {
		if err := step("Checking pull request", func() error {
			client, err := prpkg.New()
			if err != nil {
				return nil
			}
			snap.OpenPR, _ = client.ViewCurrent()
			return nil
		}); err != nil {
			return nil, err
		}
	}

	snap.NextSteps = buildNextSteps(overview, snap.OpenPR, dockerOverview, cfgErr == nil)

	if prog != nil {
		prog.Success("Ready")
	}
	return snap, nil
}

func hasGH() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

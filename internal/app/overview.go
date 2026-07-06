package app

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/laerciocrestani/gitia/internal/config"
	gitpkg "github.com/laerciocrestani/gitia/internal/git"
	prpkg "github.com/laerciocrestani/gitia/internal/pr"
	"github.com/laerciocrestani/gitia/internal/ui"
)

func RunOverview() error {
	sess := ui.New("overview", false)
	sess.Header()

	repo, err := gitpkg.New()
	if err != nil {
		return err
	}
	if err := repo.IsRepo(); err != nil {
		return fmt.Errorf("diretório atual não é um repositório git")
	}

	baseBranch := "main"
	if cfg, err := config.Load(); err == nil {
		baseBranch = cfg.BaseBranch
	}

	var overview *gitpkg.Overview
	var openPR *prpkg.PRView

	if err := sess.Step("Reading repository", func() error {
		var err error
		overview, err = repo.Overview(baseBranch)
		return err
	}); err != nil {
		return err
	}

	if hasGH() {
		_ = sess.Step("Checking pull request", func() error {
			client, err := prpkg.New()
			if err != nil {
				return nil
			}
			openPR, err = client.ViewCurrent()
			return nil
		})
	}

	fmt.Println()
	printOverview(sess, overview, openPR)
	printGitiaConfig(sess)
	printSuggestions(sess, overview, openPR)
	fmt.Println()
	sess.Success("Repository overview 📋")
	return nil
}

func printOverview(sess *ui.Session, o *gitpkg.Overview, pr *prpkg.PRView) {
	sess.Section("Repository")
	sess.KV("Path", o.Root)
	if o.RemoteURL != "" {
		sess.KV("Remote", o.RemoteURL)
	} else {
		sess.KV("Remote", "none")
	}

	sess.Section("Branch")
	if o.Detached {
		sess.KV("HEAD", "detached")
	} else {
		sess.KV("Current", o.Branch)
	}
	if o.Upstream != "" {
		sess.KV("Tracking", o.Upstream)
		sess.KV("Sync", syncLabel(o.Ahead, o.Behind))
	} else if !o.Detached {
		sess.KV("Tracking", "no upstream")
	}
	if o.CommitsAheadOfBase > 0 {
		sess.KV("vs "+o.BaseBranch, fmt.Sprintf("%d commit(s) ahead", o.CommitsAheadOfBase))
	} else if !o.Detached && o.Branch != o.BaseBranch {
		sess.KV("vs "+o.BaseBranch, "up to date")
	}

	if pr != nil {
		sess.Section("Pull request")
		state := strings.ToLower(pr.State)
		if pr.IsDraft {
			state = "draft"
		}
		sess.KV("PR", fmt.Sprintf("#%d %s", pr.Number, pr.Title))
		sess.KV("State", state)
		sess.KV("URL", pr.URL)
	}

	sess.Section("Working tree")
	sess.KV("Staged", fmt.Sprintf("%d file(s)", o.Staged))
	sess.KV("Modified", fmt.Sprintf("%d file(s)", o.Modified))
	sess.KV("Untracked", fmt.Sprintf("%d file(s)", o.Untracked))
	if !o.IsDirty() {
		sess.KV("State", "clean ✓")
	}

	if len(o.FileChanges) > 0 {
		sess.Section("Changed files")
		limit := len(o.FileChanges)
		if limit > 12 {
			limit = 12
		}
		for _, f := range o.FileChanges[:limit] {
			sess.FileChange(f.Path, f.Status, f.StatsLabel())
		}
		if len(o.FileChanges) > 12 {
			sess.Detail(fmt.Sprintf("… +%d more file(s)", len(o.FileChanges)-12))
		}
	}

	if len(o.Stashes) > 0 {
		sess.Section("Stash")
		sess.KV("Entries", fmt.Sprintf("%d saved", len(o.Stashes)))
		limit := len(o.Stashes)
		if limit > 5 {
			limit = 5
		}
		for _, stash := range o.Stashes[:limit] {
			label := stash.Ref
			if stash.Branch != "" {
				label += " on " + stash.Branch
			}
			if stash.Message != "" {
				label += ": " + stash.Message
			}
			sess.Bullet(label)
		}
		if len(o.Stashes) > 5 {
			sess.Detail(fmt.Sprintf("… +%d more stash(es)", len(o.Stashes)-5))
		}
	}

	if len(o.Branches) > 0 {
		sess.Section("Branches")
		limit := len(o.Branches)
		if limit > 8 {
			limit = 8
		}
		for _, b := range o.Branches[:limit] {
			sess.BranchLine(b.Name, b.Current, b.Upstream, b.Ahead, b.Behind)
		}
		if len(o.Branches) > 8 {
			sess.Detail(fmt.Sprintf("… +%d more", len(o.Branches)-8))
		}
	}

	if len(o.RecentCommits) > 0 {
		sess.Section("Recent commits")
		for _, line := range o.RecentCommits {
			sess.Bullet(line)
		}
	}
}

func printGitiaConfig(sess *ui.Session) {
	sess.Section("Gitia")
	cfg, err := config.Load()
	if err != nil {
		sess.KV("Config", "not configured — run: gitia config")
		return
	}
	sess.KV("Provider", string(cfg.Provider))
	sess.KV("Model", cfg.Model)
	sess.KV("Language", cfg.Language)
	sess.KV("Base branch", cfg.BaseBranch)
	sess.KV("API key", config.MaskAPIKey(cfg.APIKey))
}

func printSuggestions(sess *ui.Session, o *gitpkg.Overview, pr *prpkg.PRView) {
	var tips []string

	if _, err := config.Load(); err != nil {
		tips = append(tips, "gitia config")
	}
	if o.IsDirty() {
		tips = append(tips, "gitia commit")
	}
	if o.Ahead > 0 || (o.IsDirty() && o.Upstream != "") {
		tips = append(tips, "gitia push")
	}
	if pr == nil && o.CommitsAheadOfBase > 0 && !o.IsDirty() {
		tips = append(tips, "gitia pr")
	}
	if pr != nil {
		tips = append(tips, "gh pr view --web")
	}
	if len(o.Stashes) > 0 {
		tips = append(tips, "git stash pop")
		tips = append(tips, "git stash list")
	}
	if o.Behind > 0 {
		tips = append(tips, "git pull")
	}
	if len(tips) == 0 && !o.IsDirty() {
		tips = append(tips, "working tree clean — nothing to do")
	}

	sess.Section("Suggested next steps")
	for _, tip := range tips {
		if strings.Contains(tip, " ") && !strings.HasPrefix(tip, "gitia") && !strings.HasPrefix(tip, "git ") && !strings.HasPrefix(tip, "gh ") {
			sess.Bullet(tip)
		} else {
			sess.CommandHint(tip)
		}
	}

	if hasGH() {
		return
	}
	sess.Detail("install gh for PR info — https://cli.github.com/")
}

func syncLabel(ahead, behind int) string {
	switch {
	case ahead > 0 && behind > 0:
		return fmt.Sprintf("↑%d ahead · ↓%d behind", ahead, behind)
	case ahead > 0:
		return fmt.Sprintf("↑%d ahead of remote", ahead)
	case behind > 0:
		return fmt.Sprintf("↓%d behind remote", behind)
	default:
		return "in sync with remote"
	}
}

func hasGH() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

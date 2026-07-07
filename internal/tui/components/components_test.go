package components_test

import (
	"strings"
	"testing"

	"github.com/laerciocrestani/gitai/internal/app"
	gitpkg "github.com/laerciocrestani/gitai/internal/git"
	"github.com/laerciocrestani/gitai/internal/tui/components"
)

func TestRenderFooter(t *testing.T) {
	snap := &app.WorkspaceSnapshot{
		Overview: &gitpkg.Overview{Modified: 1},
	}
	items := components.DefaultFooterItems(snap)
	out := components.RenderFooter(items, 80)
	if !strings.Contains(out, "Commit") || !strings.Contains(out, "Quit") {
		t.Fatalf("footer missing shortcuts: %q", out)
	}
}

func TestRenderGitGraph(t *testing.T) {
	o := &gitpkg.Overview{
		Branch:             "feature/x",
		BaseBranch:         "main",
		CommitsAheadOfBase: 3,
	}
	out := components.RenderGitGraph(o, 60)
	if !strings.Contains(out, "Git Graph") || !strings.Contains(out, "HEAD") {
		t.Fatalf("git graph incomplete: %q", out)
	}
}

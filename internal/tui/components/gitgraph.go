package components

import (
	"fmt"
	"strings"

	gitpkg "github.com/laerciocrestani/openbench/internal/git"
	"github.com/laerciocrestani/openbench/internal/tui/theme"
)

// RenderGitGraph renders a compact branch graph from overview data.
func RenderGitGraph(o *gitpkg.Overview, width int) string {
	if o == nil {
		return ""
	}

	base := o.BaseBranch
	if base == "" {
		base = "main"
	}
	branch := o.Branch
	if o.Detached {
		branch = "HEAD"
	}

	commits := o.CommitsAheadOfBase
	if commits <= 0 && o.Ahead > 0 {
		commits = o.Ahead
	}
	if commits <= 0 {
		commits = 1
	}
	if commits > 6 {
		commits = 6
	}

	baseLabel := truncate(base, 12)
	branchLabel := truncate(branch, 16)

	var lines []string
	onBase := !o.Detached && o.Branch == o.BaseBranch

	if branch == base || onBase {
		dots := strings.Repeat("─", min(commits*2, 12)) + "● HEAD"
		lines = append(lines, fmt.Sprintf("%s %s", padRight(baseLabel, 14), dots))
	} else {
		baseDots := strings.Repeat("─", 12) + "●"
		lines = append(lines, fmt.Sprintf("%s %s", padRight(baseLabel, 14), baseDots))
		lines = append(lines, fmt.Sprintf("%s%s╲", strings.Repeat(" ", 14), strings.Repeat(" ", len(baseLabel))))
		branchDots := branchCommits(commits)
		lines = append(lines, fmt.Sprintf("%s %s %s", padRight(branchLabel, 14), branchDots, theme.S.Current.Render("HEAD")))
	}

	body := strings.Join(lines, "\n")
	return RenderPanel("Git Graph", body, width)
}

func branchCommits(n int) string {
	if n <= 1 {
		return "●"
	}
	parts := make([]string, n)
	for i := range parts {
		parts[i] = "●"
	}
	return strings.Join(parts, "──")
}

func padRight(s string, n int) string {
	for len(s) < n {
		s += " "
	}
	return s
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

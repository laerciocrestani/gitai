package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	gitpkg "github.com/laerciocrestani/gitai/internal/git"
	"github.com/laerciocrestani/gitai/internal/tui/theme"
)

// RenderFileTable renders changed files as an aligned table.
func RenderFileTable(changes []gitpkg.FileChange, width, maxRows int) string {
	if len(changes) == 0 {
		return ""
	}

	if maxRows <= 0 {
		maxRows = 12
	}
	limit := len(changes)
	if limit > maxRows {
		limit = maxRows
	}

	inner := width - 4
	if inner < 40 {
		inner = 40
	}

	pathWidth := inner - 24
	if pathWidth < 20 {
		pathWidth = 20
	}

	header := fmt.Sprintf("%-4s %-*s %6s %6s", "TYPE", pathWidth, "FILE", "+", "-")
	var rows []string
	rows = append(rows, theme.S.Hint.Render(header))

	for _, f := range changes[:limit] {
		tag := statusTag(f.Status)
		path := truncate(f.Path, pathWidth)
		plus := fmt.Sprintf("%d", f.Insertions)
		minus := fmt.Sprintf("%d", f.Deletions)
		row := fmt.Sprintf("%-4s %-*s %6s %6s", tag, pathWidth, path, plus, minus)
		rows = append(rows, fileRowStyle(f.Status).Render(row))
	}

	footer := fmt.Sprintf("Total: %d files", len(changes))
	if len(changes) > limit {
		footer += fmt.Sprintf(" (showing %d)", limit)
	}
	rows = append(rows, theme.S.Hint.Render(footer))

	body := strings.Join(rows, "\n")
	return RenderPanel("Changed Files", body, width)
}

func fileRowStyle(status string) lipgloss.Style {
	switch status {
	case "untracked":
		return theme.S.Untracked
	case "deleted":
		return theme.S.Error
	case "new", "staged":
		return theme.S.New
	case "modified", "staged+modified":
		return theme.S.Modified
	default:
		return theme.S.Hint
	}
}

func statusTag(status string) string {
	switch status {
	case "untracked":
		return "?"
	case "deleted":
		return "D"
	case "new", "staged":
		return "A"
	case "modified", "staged+modified":
		return "M"
	default:
		return "·"
	}
}

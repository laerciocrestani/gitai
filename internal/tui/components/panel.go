package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/laerciocrestani/gitai/internal/tui/theme"
)

// RenderPanel renders a titled panel with optional body content.
func RenderPanel(title, body string, width int) string {
	if width < 20 {
		width = 78
	}
	inner := width - 4
	titleLine := "╭ " + title
	dashLen := inner - runewidth.StringWidth(title) - 1
	if dashLen < 1 {
		dashLen = 1
	}
	titleLine += " " + strings.Repeat("─", dashLen) + "╮"

	var lines []string
	lines = append(lines, theme.S.PanelTitle.Render(titleLine))

	if body == "" {
		lines = append(lines, theme.S.Panel.Render(boxEmpty(inner)))
	} else {
		for _, line := range strings.Split(strings.TrimSuffix(body, "\n"), "\n") {
			lines = append(lines, theme.S.Panel.Render(boxContent(line, inner)))
		}
	}
	lines = append(lines, theme.S.Panel.Render(boxBottom(inner)))
	return strings.Join(lines, "\n") + "\n"
}

func boxEmpty(inner int) string {
	padding := inner
	if padding < 0 {
		padding = 0
	}
	return "│" + strings.Repeat(" ", padding) + "│"
}

func boxContent(content string, inner int) string {
	w := runewidth.StringWidth(content)
	if w > inner {
		content = truncate(content, inner)
		w = runewidth.StringWidth(content)
	}
	pad := inner - w
	if pad < 0 {
		pad = 0
	}
	return "│" + content + strings.Repeat(" ", pad) + "│"
}

func boxBottom(inner int) string {
	return "╰" + strings.Repeat("─", inner) + "╯"
}

func truncate(s string, max int) string {
	return runewidth.Truncate(s, max, "…")
}

// RenderDivider renders a horizontal divider spanning the given width.
func RenderDivider(width int) string {
	if width < 4 {
		width = 78
	}
	return theme.S.Hint.Render("├" + strings.Repeat("─", width-2) + "┤") + "\n"
}

// PadLine pads content to the given display width.
func PadLine(left, right string, width int) string {
	gap := width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	return left + strings.Repeat(" ", gap) + right
}

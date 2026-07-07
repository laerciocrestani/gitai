package components

import (
	"strings"

	"github.com/laerciocrestani/gitai/internal/tui/theme"
)

// RenderLoading renders a loading state with an animated progress bar.
func RenderLoading(message string, tick int, width int) string {
	if width < 20 {
		width = 78
	}
	barWidth := width - 8
	if barWidth > 40 {
		barWidth = 40
	}
	if barWidth < 10 {
		barWidth = 10
	}

	filled := tick % (barWidth + 1)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	var lines []string
	lines = append(lines, theme.S.Hint.Render(message))
	lines = append(lines, theme.S.Info.Render(bar))

	body := strings.Join(lines, "\n")
	return RenderPanel("Loading", body, width)
}

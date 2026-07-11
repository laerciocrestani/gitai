package components

import (
	"fmt"
	"strings"

	"github.com/laerciocrestani/openbench/internal/app"
	"github.com/laerciocrestani/openbench/internal/tui/theme"
)

// RenderAIPanel renders the AI engine status panel.
func RenderAIPanel(snap *app.WorkspaceSnapshot, width int) string {
	if snap == nil {
		return ""
	}

	var parts []string

	if snap.ConfigErr != nil {
		parts = append(parts, theme.S.Warn.Render("Status: "+snap.ConfigErr.Error()))
	} else if snap.Config != nil {
		model := snap.Config.Model
		status := "Ready"
		if snap.Config.APIKey == "" {
			status = "Setup required"
		}
		label := app.FormatProviderName(snap.Config.Provider)
		if model != "" {
			label += " " + model
		}
		parts = append(parts, fmt.Sprintf("%s │ %s", label, theme.S.Success.Render(status)))

		if ctx := app.ModelContextWindow(model); ctx != "" {
			parts = append(parts, fmt.Sprintf("Context %s", ctx))
		}
		if cost := app.EstimateAICost(snap); cost != "" {
			parts = append(parts, fmt.Sprintf("Cost %s", cost))
		}
	} else {
		parts = append(parts, theme.S.Warn.Render("Status: not configured"))
	}

	body := strings.Join(parts, " │ ")
	return RenderPanel("AI Engine", body, width)
}

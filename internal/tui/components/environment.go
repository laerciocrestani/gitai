package components

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/laerciocrestani/openbench/internal/app"
	dockerpkg "github.com/laerciocrestani/openbench/internal/docker"
	"github.com/laerciocrestani/openbench/internal/tui/theme"
)

// RenderEnvironmentPanel renders Docker/Compose status at the top of the dashboard.
func RenderEnvironmentPanel(snap *app.WorkspaceSnapshot, width int) string {
	if snap == nil || snap.Docker == nil {
		return ""
	}
	ov := snap.Docker
	var lines []string

	daemon := "missing"
	if ov.Available {
		if ov.DaemonRunning {
			daemon = theme.S.Success.Render("running")
		} else {
			daemon = theme.S.Error.Render("stopped")
		}
	} else {
		daemon = theme.S.Disabled.Render("n/a")
	}
	lines = append(lines, fmt.Sprintf("Docker  %s", daemon))

	if ov.ComposeFile != "" {
		composeName := filepath.Base(ov.ComposeFile)
		lines = append(lines, fmt.Sprintf("Compose  %s  (%s)", composeName, ov.ProjectName))
	} else if ov.Available && ov.DaemonRunning {
		lines = append(lines, theme.S.Hint.Render("Compose  not found"))
	}

	if ov.Error != "" && len(ov.Containers) == 0 {
		lines = append(lines, theme.S.Warn.Render(ov.Error))
	}

	if len(ov.Containers) == 0 {
		note := app.FormatDockerNote(ov)
		if note != "" {
			lines = append(lines, theme.S.Hint.Render(note))
		}
	} else {
		maxRows := 6
		for i, c := range ov.Containers {
			if i >= maxRows {
				lines = append(lines, theme.S.Hint.Render(fmt.Sprintf("… +%d more", len(ov.Containers)-maxRows)))
				break
			}
			lines = append(lines, formatContainerLine(c))
		}
	}

	body := strings.Join(lines, "\n")
	return RenderPanel("Environment", body, width)
}

func formatContainerLine(c dockerpkg.ContainerSummary) string {
	stateStyle := theme.S.Hint
	switch strings.ToLower(c.State) {
	case "running":
		stateStyle = theme.S.Success
	case "exited", "dead":
		stateStyle = theme.S.Error
	}
	line := fmt.Sprintf("%-10s %s", c.Service, stateStyle.Render(c.State))
	if c.Ports != "" {
		line += "  " + theme.S.Info.Render(c.Ports)
	}
	if c.Health != "" {
		line += "  " + theme.S.Hint.Render("("+c.Health+")")
	}
	return line
}

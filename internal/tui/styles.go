package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/laerciocrestani/openbench/internal/tui/theme"
)

var (
	styleTitle     lipgloss.Style
	styleHeader    lipgloss.Style
	styleSection   lipgloss.Style
	styleCurrent   lipgloss.Style
	styleHint      lipgloss.Style
	styleStatusBar lipgloss.Style
	styleError     lipgloss.Style
	styleKey       lipgloss.Style
	styleModified  lipgloss.Style
	styleNew       lipgloss.Style
	styleUntracked lipgloss.Style
	styleYellow    lipgloss.Style
	styleWarn      lipgloss.Style
	stylePanel     lipgloss.Style
	styleSuccess   lipgloss.Style
	styleInfo      lipgloss.Style
)

func init() {
	initTheme()
}

func initTheme() {
	theme.Init()
	syncStyles()
}

func syncStyles() {
	styleTitle = theme.S.Title
	styleHeader = theme.S.Header
	styleSection = theme.S.Section
	styleCurrent = theme.S.Current
	styleHint = theme.S.Hint
	styleStatusBar = theme.S.StatusBar
	styleError = theme.S.Error
	styleKey = theme.S.Key
	styleModified = theme.S.Modified
	styleNew = theme.S.New
	styleUntracked = theme.S.Untracked
	styleYellow = theme.S.Yellow
	styleWarn = theme.S.Warn
	stylePanel = theme.S.Panel
	styleSuccess = theme.S.Success
	styleInfo = theme.S.Info
}

func themePlain() bool {
	return theme.Plain()
}

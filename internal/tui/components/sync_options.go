package components

import (
	"fmt"
	"strings"

	"github.com/laerciocrestani/gitai/internal/tui/theme"
	"github.com/laerciocrestani/gitai/internal/ui"
)

// SyncMode identifies a sync execution preset.
type SyncMode int

const (
	SyncModeStandard SyncMode = iota
	SyncModePruneRemote
	SyncModePruneFull
)

// SyncModeOption describes one sync preset with CLI flag and explanation.
type SyncModeOption struct {
	Mode        SyncMode
	Label       string
	Flag        string
	Summary     string
	Description string
	Prune       bool
	PruneRemote bool
}

// SyncModeCatalog returns all sync presets in display order.
func SyncModeCatalog() []SyncModeOption {
	return []SyncModeOption{
		{
			Mode:        SyncModeStandard,
			Label:       "Sync padrão",
			Flag:        "(nenhuma)",
			Summary:     "Fetch + pull da branch base",
			Description: "Atualiza refs remotas (fetch --prune) e faz fast-forward da branch base com origin. Não remove branches mergeadas.",
			Prune:       false,
			PruneRemote: false,
		},
		{
			Mode:        SyncModePruneRemote,
			Label:       "Sync + prune remoto",
			Flag:        "--prune-remote",
			Summary:     "Sync + limpa branches no GitHub",
			Description: "Após o sync, remove branches remotas já mergeadas na base (git push origin --delete). Mantém branches locais.",
			Prune:       false,
			PruneRemote: true,
		},
		{
			Mode:        SyncModePruneFull,
			Label:       "Sync + prune completo",
			Flag:        "--prune",
			Summary:     "Sync + limpa local e remoto",
			Description: "Após o sync, remove branches locais e remotas já mergeadas na base. Branches divergentes pedem confirmação antes do -D.",
			Prune:       true,
			PruneRemote: false,
		},
	}
}

// ToAppOptions maps the selected preset to app.SyncOptions fields.
func (o SyncModeOption) ToAppOptions(base string) (prune, pruneRemote bool, resolvedBase string) {
	return o.Prune, o.PruneRemote || o.Prune, base
}

// RenderSyncOptionsPanel renders the sync mode picker with a detail table.
func RenderSyncOptionsPanel(cursor int, modes []SyncModeOption, base string, dirty bool, width int) string {
	inner := ui.ContentInner(width)
	var lines []string

	if dirty {
		lines = append(lines, theme.S.Warn.Render("  ⚠ Working tree suja — commit ou stash antes de sincronizar"))
		lines = append(lines, "")
	}

	for i, mode := range modes {
		marker := "  "
		if i == cursor {
			marker = "> "
		}
		flag := mode.Flag
		if flag == "" {
			flag = theme.S.Hint.Render("(padrão)")
		} else {
			flag = theme.S.Key.Render(flag)
		}
		label := mode.Label + "  " + flag
		if i == cursor {
			lines = append(lines, theme.S.Current.Render(marker+label))
		} else {
			lines = append(lines, theme.S.Hint.Render(marker+label))
		}
	}

	lines = append(lines, "")
	lines = append(lines, theme.S.Hint.Render("  Base: "+base))
	lines = append(lines, "")

	selected := modes[cursor]
	lines = append(lines, renderSyncDetailTable(selected, base, inner))

	body := strings.Join(lines, "\n")
	return RenderPanel("Sync · Opções", body, width)
}

func renderSyncDetailTable(mode SyncModeOption, base string, inner int) string {
	const (
		colW = 14
	)

	lines := []string{
		theme.S.Hint.Render(fmt.Sprintf("  %-*s %s", colW, "Opção", mode.Label)),
		theme.S.Hint.Render(fmt.Sprintf("  %-*s %s", colW, "Flag", mode.Flag)),
		theme.S.Hint.Render(fmt.Sprintf("  %-*s %s", colW, "Resumo", truncatePlain(mode.Summary, inner-colW-2))),
		"",
		theme.S.Hint.Render("  O que faz"),
		"  " + wrapPlain(mode.Description, inner-2),
		"",
		theme.S.Hint.Render("  Comandos"),
	}

	for _, cmd := range syncCommandPreview(mode, base) {
		lines = append(lines, theme.S.Hint.Render("  · "+cmd))
	}

	return strings.Join(lines, "\n")
}

func syncCommandPreview(mode SyncModeOption, base string) []string {
	if base == "" {
		base = "main"
	}
	cmds := []string{
		"git fetch origin --prune",
		"git checkout " + base,
		"git pull --ff-only origin " + base,
	}
	if mode.Prune || mode.PruneRemote {
		cmds = append(cmds, "git branch --merged " + base + " …")
	}
	if mode.Prune {
		cmds = append(cmds, "git branch -d/-D <merged-local> …")
	}
	if mode.Prune || mode.PruneRemote {
		cmds = append(cmds, "git push origin --delete <merged-remote> …")
	}
	return cmds
}

// RenderSyncBaseEditor renders the base branch edit step.
func RenderSyncBaseEditor(baseField string, width int) string {
	body := theme.S.Hint.Render("  Branch base para pull e prune:\n\n  ") + baseField
	return RenderPanel("Sync · Base branch", body, width)
}

func wrapPlain(text string, width int) string {
	if width < 20 {
		return text
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}
	var lines []string
	var current strings.Builder
	for _, word := range words {
		add := word
		if current.Len() > 0 {
			add = " " + word
		}
		if current.Len()+len(add) > width && current.Len() > 0 {
			lines = append(lines, current.String())
			current.Reset()
			current.WriteString(word)
			continue
		}
		current.WriteString(add)
	}
	if current.Len() > 0 {
		lines = append(lines, current.String())
	}
	return strings.Join(lines, "\n")
}

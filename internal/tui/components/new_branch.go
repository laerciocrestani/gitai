package components

import (
	"fmt"
	"strings"

	"github.com/laerciocrestani/gitai/internal/tui/theme"
	"github.com/laerciocrestani/gitai/internal/ui"
)

// NewBranchStep identifies the wizard step for creating a branch.
type NewBranchStep int

const (
	NewBranchStepFrom NewBranchStep = iota
	NewBranchStepTemplate
	NewBranchStepName
)

// RenderNewBranchFromPanel renders the source branch picker.
func RenderNewBranchFromPanel(cursor, total int, body string, width int) string {
	title := "New Branch · From"
	if total > 0 {
		title += fmt.Sprintf("  %d/%d", cursor+1, total)
	}
	if strings.TrimSpace(body) == "" {
		body = theme.S.Hint.Render("  (nenhuma branch local)")
	}
	return RenderPanel(title, body, width)
}

// RenderNewBranchTemplateListBody renders only the scrollable template rows.
func RenderNewBranchTemplateListBody(cursor int, items []NewBranchTemplateItem) string {
	var lines []string
	selectableIdx := 0
	for _, item := range items {
		if item.Separator {
			lines = append(lines, theme.S.Hint.Render("  ────────────────"))
			continue
		}
		line := "  " + item.Template.Label()
		if selectableIdx == cursor {
			line = theme.S.Current.Render("> " + strings.TrimPrefix(line, "  "))
		} else {
			line = theme.S.Hint.Render(line)
		}
		lines = append(lines, line)
		selectableIdx++
	}
	return strings.Join(lines, "\n")
}

// RenderNewBranchTemplatePanel renders the template list viewport and detail table.
func RenderNewBranchTemplatePanel(cursor, selectable int, items []NewBranchTemplateItem, listBody string, selected NewBranchTemplate, width int) string {
	inner := ui.ContentInner(width)
	var lines []string
	if strings.TrimSpace(listBody) != "" {
		lines = append(lines, listBody)
	}
	lines = append(lines, "")
	lines = append(lines, renderTemplateTable(selected, inner))

	body := strings.Join(lines, "\n")
	title := "New Branch · Template"
	if selectable > 0 {
		title += fmt.Sprintf("  %d/%d", cursor+1, selectable)
	}
	return RenderPanel(title, body, width)
}

func renderTemplateTable(t NewBranchTemplate, inner int) string {
	const (
		prefixW = 16
		usageW  = 34
	)

	header := fmt.Sprintf("%-*s %-*s %s", prefixW, "Prefixo", usageW, "Uso", "Exemplo")
	lines := []string{theme.S.Hint.Render("  "+header)}

	if t.Icon == "" && !t.Other {
		lines = append(lines, theme.S.Hint.Render("  Selecione um template"))
		return strings.Join(lines, "\n")
	}

	prefix := truncatePlain(t.PrefixColumn(), prefixW)
	usage := truncatePlain(t.Usage, usageW)
	example := t.Example
	if t.Other {
		example = "(livre)"
	}
	row := fmt.Sprintf("%-*s %-*s %s", prefixW, prefix, usageW, usage, example)
	lines = append(lines, "  "+row)
	return strings.Join(lines, "\n")
}

// RenderNewBranchNamePanel renders the branch name input step.
func RenderNewBranchNamePanel(from string, template NewBranchTemplate, nameField string, width int) string {
	var lines []string

	lines = append(lines, theme.S.Hint.Render("  From: "+from))
	if template.Other {
		lines = append(lines, theme.S.Hint.Render("  Template: Outro (nome livre)"))
	} else if template.Prefix != "" {
		lines = append(lines, theme.S.Hint.Render("  Template: "+template.Label()))
	}
	lines = append(lines, "")
	lines = append(lines, theme.S.Hint.Render("  Nome da branch:"))
	lines = append(lines, "  "+nameField)

	preview := strings.TrimSpace(stripANSI(nameField))
	if preview == "" {
		preview = "(digite o nome)"
	}
	lines = append(lines, "")
	lines = append(lines, theme.S.Hint.Render("  Preview: "+preview))

	body := strings.Join(lines, "\n")
	return RenderPanel("New Branch · Name", body, width)
}

func truncatePlain(s string, max int) string {
	if ui.DisplayWidth(s) <= max {
		return s
	}
	runes := []rune(s)
	for len(runes) > 0 && ui.DisplayWidth(string(runes)) > max-1 {
		runes = runes[:len(runes)-1]
	}
	return string(runes) + "…"
}

func stripANSI(s string) string {
	var b strings.Builder
	esc := false
	for _, r := range s {
		if r == '\x1b' {
			esc = true
			continue
		}
		if esc {
			if r == 'm' {
				esc = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

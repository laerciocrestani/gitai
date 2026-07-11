package tui

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/laerciocrestani/openbench/internal/app"
	"github.com/laerciocrestani/openbench/internal/tui/components"
)

type doctorModel struct {
	viewport viewport.Model
	content  string
	ready    bool
	err      error
	explain  bool
}

func newDoctorModel() doctorModel {
	return doctorModel{}
}

func loadDoctorCmd(explain bool) tea.Cmd {
	return func() tea.Msg {
		prog := NewActionProgress()
		prog.Reset()
		report, err := app.RunDoctor(context.Background(), app.DoctorOptions{
			Explain:  explain,
			Progress: prog,
		})
		if err != nil {
			return doctorLoadedMsg{err: err}
		}
		return doctorLoadedMsg{
			content: app.FormatDoctorContent(report),
			explain: explain,
		}
	}
}

type doctorLoadedMsg struct {
	content string
	err     error
	explain bool
}

func (m *doctorModel) SetSize(width, height int) {
	headerRows := 4
	footerRows := 2
	vh := height - headerRows - footerRows
	if vh < 3 {
		vh = 3
	}
	if !m.ready {
		m.viewport = viewport.New(width, vh)
		m.viewport.SetContent(m.content)
		m.ready = true
	} else {
		m.viewport.Width = width
		m.viewport.Height = vh
	}
}

func (m *doctorModel) Load(msg doctorLoadedMsg) {
	m.content = msg.content
	m.err = msg.err
	m.explain = msg.explain
	m.ready = false
}

func (m doctorModel) Update(msg tea.Msg) (doctorModel, tea.Cmd) {
	if m.err != nil {
		return m, nil
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m doctorModel) View(tick int) string {
	if m.err != nil {
		return styleError.Render("  ✗ " + m.err.Error())
	}

	var b strings.Builder
	title := "Repository Health"
	if m.explain {
		title += " · AI"
	}
	b.WriteString(styleSection.Render(title))
	b.WriteString("\n")
	b.WriteString(styleHint.Render("  Panorama do que está sendo desenvolvido"))
	b.WriteString("\n\n")
	if !m.ready {
		b.WriteString(components.RenderSpinnerLine("Analyzing", tick))
		return b.String()
	}
	b.WriteString(m.viewport.View())
	return b.String()
}

func doctorHelpLine(explain bool) string {
	line := styleKey.Render("e") + " explain with AI  " +
		styleKey.Render("r") + " refresh  " +
		styleKey.Render("esc") + " back"
	if explain {
		line = styleKey.Render("r") + " refresh  " +
			styleKey.Render("esc") + " back"
	}
	return line
}

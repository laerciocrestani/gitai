package ui

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const boxWidth = 44

type Session struct {
	command string
	dryRun  bool
	enabled bool
}

func New(command string, dryRun bool) *Session {
	return &Session{
		command: command,
		dryRun:  dryRun,
		enabled: colorsEnabled(),
	}
}

func colorsEnabled() bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("GITIA_NO_UI") != "" {
		return false
	}
	if os.Getenv("CI") != "" {
		return false
	}
	fi, err := os.Stderr.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func (s *Session) Header() {
	subtitle := strings.ToLower(s.command)
	if s.dryRun {
		subtitle += " · dry-run"
	}
	if !s.enabled {
		fmt.Fprintf(os.Stderr, "\n🤖 Gitia %s · %s\n\n", Version, subtitle)
		return
	}
	s.printBox(
		fmt.Sprintf("🤖 Gitia %s · %s", Version, subtitle),
		"",
	)
	fmt.Fprintln(os.Stderr)
}

func (s *Session) Step(label string, fn func() error) error {
	if !s.enabled {
		return fn()
	}

	stop := s.spinner(label)
	err := fn()
	stop()

	if err != nil {
		s.failLine(label)
		return err
	}
	s.doneLine(label)
	return nil
}

func (s *Session) Info(label string) {
	if !s.enabled {
		fmt.Fprintln(os.Stderr, label)
		return
	}
	fmt.Fprintln(os.Stderr, s.paint("  • "+label, dim))
}

func (s *Session) Success(message string) {
	if !s.enabled {
		fmt.Fprintf(os.Stderr, "\n%s\n\n", message)
		return
	}
	fmt.Fprintln(os.Stderr)
	s.printBox(message, "")
	fmt.Fprintln(os.Stderr)
}

func (s *Session) Detail(message string) {
	if !s.enabled {
		fmt.Fprintln(os.Stderr, message)
		return
	}
	fmt.Fprintln(os.Stderr, s.paint("  → "+message, cyan))
}

func (s *Session) Warn(message string) {
	if !s.enabled {
		fmt.Fprintf(os.Stderr, "! %s\n", message)
		return
	}
	fmt.Fprintln(os.Stderr, s.paint("  ! "+message, yellow))
}

func (s *Session) Prompt(label string) {
	if !s.enabled {
		fmt.Print(label)
		return
	}
	fmt.Fprint(os.Stderr, s.paint("  ? "+label, magenta))
}

func (s *Session) UsageBlock(lines []string) {
	if !s.enabled {
		for _, line := range lines {
			fmt.Println(line)
		}
		return
	}
	fmt.Println()
	s.printBoxLines(append([]string{s.paint("📊 Uso de IA", bold+cyan)}, lines...))
	fmt.Println()
}

func (s *Session) Section(title string) {
	fmt.Println()
	if s.enabled {
		fmt.Println(s.paint("▸ "+title, bold+cyan))
	} else {
		fmt.Println("▸ " + title)
	}
}

func (s *Session) KV(key, value string) {
	if s.enabled {
		fmt.Printf("  %s %s\n", s.paint(key+":", dim), value)
		return
	}
	fmt.Printf("  %s: %s\n", key, value)
}

func (s *Session) Bullet(text string) {
	if s.enabled {
		fmt.Println("  " + s.paint("•", dim) + " " + text)
		return
	}
	fmt.Println("  • " + text)
}

func (s *Session) BranchLine(name string, current bool, upstream string, ahead, behind int) {
	marker := " "
	if current {
		marker = s.paint("*", green)
	}

	line := fmt.Sprintf("%s %s", marker, name)
	if upstream != "" {
		line += s.paint(" → "+upstream, dim)
	}
	if ahead > 0 || behind > 0 {
		line += s.paint(fmt.Sprintf(" (↑%d ↓%d)", ahead, behind), yellow)
	}
	fmt.Println("  " + line)
}

func (s *Session) CommandHint(cmd string) {
	if s.enabled {
		fmt.Println("  " + s.paint("→", cyan) + " " + s.paint(cmd, bold+magenta))
		return
	}
	fmt.Println("  → " + cmd)
}

func (s *Session) FileChange(path, status, stats string) {
	statusLabel := s.paint("["+status+"]", fileStatusColor(status))
	line := "  " + statusLabel + " " + path
	if stats != "" {
		line += " " + s.paint(stats, green)
	}
	fmt.Println(line)
}

func fileStatusColor(status string) string {
	switch status {
	case "untracked":
		return yellow
	case "deleted":
		return red
	case "new", "staged":
		return green
	case "modified", "staged+modified":
		return magenta
	default:
		return cyan
	}
}

func (s *Session) printBox(lines ...string) {
	s.printBoxLines(lines)
}

func (s *Session) printBoxLines(lines []string) {
	top := "╭" + strings.Repeat("─", boxWidth) + "╮"
	bottom := "╰" + strings.Repeat("─", boxWidth) + "╯"
	fmt.Fprintln(os.Stderr, s.paint(top, cyan))
	for _, line := range lines {
		if line == "" {
			fmt.Fprintln(os.Stderr, s.paint("│"+strings.Repeat(" ", boxWidth)+"│", cyan))
			continue
		}
		fmt.Fprintln(os.Stderr, s.boxLine(line))
	}
	fmt.Fprintln(os.Stderr, s.paint(bottom, cyan))
}

func (s *Session) boxLine(content string) string {
	inner := boxWidth - 2
	plain := stripANSI(content)
	pad := inner - visibleLen(plain)
	if pad < 0 {
		pad = 0
	}
	return s.paint("│", cyan) + " " + content + strings.Repeat(" ", pad) + " " + s.paint("│", cyan)
}

func (s *Session) spinner(label string) func() {
	if !s.enabled {
		return func() {}
	}

	done := make(chan struct{})
	var once sync.Once
	stop := func() {
		once.Do(func() { close(done) })
	}

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	go func() {
		i := 0
		ticker := time.NewTicker(90 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				fmt.Fprint(os.Stderr, "\r\033[K")
				return
			case <-ticker.C:
				frame := s.paint(frames[i%len(frames)], cyan)
				fmt.Fprintf(os.Stderr, "\r  %s %s", frame, s.paint(label+"...", yellow))
				i++
			}
		}
	}()

	return stop
}

func (s *Session) doneLine(label string) {
	fmt.Fprintf(os.Stderr, "  %s %s\n", s.paint("✓", green), s.paint(label, green))
}

func (s *Session) failLine(label string) {
	fmt.Fprintf(os.Stderr, "  %s %s\n", s.paint("✗", red), s.paint(label, red))
}

func (s *Session) paint(text, code string) string {
	if !s.enabled {
		return text
	}
	return code + text + reset
}

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	cyan    = "\033[36m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	magenta = "\033[35m"
	red     = "\033[31m"
)

func stripANSI(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	skip := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			skip = true
			continue
		}
		if skip {
			if s[i] == 'm' {
				skip = false
			}
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func visibleLen(s string) int {
	n := 0
	for _, r := range s {
		if r > 0xFFFF {
			n += 2
		} else {
			n++
		}
	}
	return n
}

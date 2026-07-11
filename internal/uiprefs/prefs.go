package uiprefs

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	defaultAutoRefreshSeconds = 5
	localConfigName           = ".openbench.yaml"
)

type filePrefs struct {
	InteractiveUI        *bool `yaml:"interactive_ui"`
	UIColor              *bool `yaml:"ui_color"`
	UIAutoRefreshSeconds *int  `yaml:"ui_auto_refresh_seconds"`
	UIWatchFiles         *bool `yaml:"ui_watch_files"`
}

// InteractiveUIEnabled reports whether `ob` without subcommand should open the TUI.
func InteractiveUIEnabled() bool {
	if os.Getenv("OB_NO_UI") != "" || os.Getenv("CI") != "" {
		return false
	}
	prefs := loadPrefs()
	if prefs.InteractiveUI == nil {
		return true
	}
	return *prefs.InteractiveUI
}

// AutoRefreshInterval returns the dashboard TUI polling interval.
func AutoRefreshInterval() time.Duration {
	secs := defaultAutoRefreshSeconds
	prefs := loadPrefs()
	if prefs.UIAutoRefreshSeconds != nil {
		secs = *prefs.UIAutoRefreshSeconds
	}
	if secs <= 0 {
		return 0
	}
	return time.Duration(secs) * time.Second
}

// WatchFilesEnabled reports whether the TUI watches filesystem changes.
func WatchFilesEnabled() bool {
	prefs := loadPrefs()
	if prefs.UIWatchFiles == nil {
		return true
	}
	return *prefs.UIWatchFiles
}

// ColorsEnabled reports whether ANSI/lipgloss colors are active.
func ColorsEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	prefs := loadPrefs()
	if prefs.UIColor == nil {
		return true
	}
	return *prefs.UIColor
}

func loadPrefs() filePrefs {
	var merged filePrefs
	for _, path := range configPaths() {
		p, ok, err := readPrefsFile(path)
		if err != nil || !ok {
			continue
		}
		if p.InteractiveUI != nil {
			merged.InteractiveUI = p.InteractiveUI
		}
		if p.UIColor != nil {
			merged.UIColor = p.UIColor
		}
		if p.UIAutoRefreshSeconds != nil {
			merged.UIAutoRefreshSeconds = p.UIAutoRefreshSeconds
		}
		if p.UIWatchFiles != nil {
			merged.UIWatchFiles = p.UIWatchFiles
		}
	}
	return merged
}

func readPrefsFile(path string) (filePrefs, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return filePrefs{}, false, nil
		}
		return filePrefs{}, false, err
	}
	var p filePrefs
	if err := yaml.Unmarshal(data, &p); err != nil {
		return filePrefs{}, false, err
	}
	return p, true, nil
}

func configPaths() []string {
	local := localConfigPath()
	if local != "" {
		if _, err := os.Stat(local); err == nil {
			return []string{local}
		}
	}
	if env := os.Getenv("OB_CONFIG"); env != "" {
		return []string{env}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	return []string{filepath.Join(home, ".config", "openbench", "config.yaml")}
}

func localConfigPath() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Join(wd, localConfigName)
}

package ui

import (
	"runtime/debug"
	"strings"
)

// CurrentVersion é a versão oficial do gitia — atualize aqui a cada release.
const CurrentVersion = "v0.1.0"

// Injetados via -ldflags no go install.
var (
	buildVersion string
	buildCommit  string
)

func Version() string {
	v := resolveVersion()
	if commit := resolveCommit(); commit != "" {
		return v + " · " + commit
	}
	return v
}

// VersionShort retorna só a versão sem sufixo de commit.
func VersionShort() string {
	return resolveVersion()
}

func resolveVersion() string {
	if v := strings.TrimSpace(buildVersion); v != "" && !isPseudoVersion(v) {
		return normalizeTag(v)
	}

	info, ok := debug.ReadBuildInfo()
	if ok {
		if v := info.Main.Version; v != "" && v != "(devel)" && !isPseudoVersion(v) {
			return normalizeTag(v)
		}
	}

	return CurrentVersion
}

func resolveCommit() string {
	if c := strings.TrimSpace(buildCommit); c != "" {
		return shortHash(c)
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	if rev := vcsSetting(info, "vcs.revision"); rev != "" {
		return shortHash(rev)
	}
	return ""
}

func SetBuildVersion(v string) {
	buildVersion = v
}

func SetBuildCommit(c string) {
	buildCommit = c
}

func isPseudoVersion(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" || v == "(devel)" {
		return false
	}
	// v0.0.0-20260706025933-a2ee046f6eb2
	if strings.Count(v, "-") >= 2 {
		return true
	}
	return strings.Contains(v, "0.0.0-")
}

func normalizeTag(v string) string {
	v = strings.TrimSpace(v)
	if v == "" || v == "(devel)" {
		return CurrentVersion
	}
	if strings.HasPrefix(v, "v") {
		return v
	}
	if strings.Contains(v, ".") && !strings.Contains(v, "/") {
		return "v" + v
	}
	return v
}

func shortHash(rev string) string {
	rev = strings.TrimSpace(rev)
	if len(rev) > 7 {
		return rev[:7]
	}
	return rev
}

func vcsSetting(info *debug.BuildInfo, key string) string {
	for _, s := range info.Settings {
		if s.Key == key {
			return s.Value
		}
	}
	return ""
}

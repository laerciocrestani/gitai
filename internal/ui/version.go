package ui

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/laerciocrestani/gitia/internal/version"
)

// Injetados via -ldflags no go install.
var (
	buildVersion string
	buildCommit  string
)

var (
	runtimeOnce sync.Once
	runtimeInfo version.Info
	runtimeOK   bool
)

func Version() string {
	if v := strings.TrimSpace(buildVersion); v != "" {
		if c := strings.TrimSpace(buildCommit); c != "" {
			return v + " · " + shortHash(c)
		}
		return v
	}

	info, ok := resolveRuntime()
	if ok {
		return info.Display()
	}
	return "v" + version.DefaultBase
}

func VersionInfo() version.Info {
	if v := strings.TrimSpace(buildVersion); v != "" {
		return version.Info{
			Version: v,
			Commit:  shortHash(buildCommit),
		}
	}
	if info, ok := resolveRuntime(); ok {
		return info
	}
	return version.Info{Version: "v" + version.DefaultBase}
}

func SetBuildVersion(v string) {
	buildVersion = v
}

func SetBuildCommit(c string) {
	buildCommit = c
}

func resolveRuntime() (version.Info, bool) {
	runtimeOnce.Do(func() {
		if root := findGitiaRoot(); root != "" {
			if info, err := version.Compute(root); err == nil {
				runtimeInfo = info
				runtimeOK = true
			}
		}
	})
	return runtimeInfo, runtimeOK
}

func findGitiaRoot() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == "github.com/laerciocrestani/gitia" && dep.Sum != "" {
				// módulo em cache — tenta cwd
				break
			}
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		modPath := filepath.Join(dir, "go.mod")
		data, err := os.ReadFile(modPath)
		if err == nil && strings.Contains(string(data), "github.com/laerciocrestani/gitia") {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func shortHash(rev string) string {
	rev = strings.TrimSpace(rev)
	if len(rev) > 7 {
		return rev[:7]
	}
	return rev
}

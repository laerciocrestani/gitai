package version

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// DefaultBase é a versão inicial quando ainda não há tags no repositório.
const DefaultBase = "0.1.0"

type Info struct {
	Version      string
	Commit       string
	Tag          string
	CommitsSince int
	Dirty        bool
	ExactTag     bool
}

func (i Info) Display() string {
	if i.ShowCommit() {
		return i.Version + " · " + i.Commit
	}
	return i.Version
}

func (i Info) ShowCommit() bool {
	return i.CommitsSince > 0 || i.Dirty
}

func (i Info) LDFlags() string {
	flags := fmt.Sprintf("-X github.com/laerciocrestani/gitia/internal/ui.buildVersion=%s", i.Version)
	if i.ShowCommit() && i.Commit != "" {
		flags += fmt.Sprintf(" -X github.com/laerciocrestani/gitia/internal/ui.buildCommit=%s", i.Commit)
	}
	return flags
}

func Compute(repoDir string) (Info, error) {
	commit, err := gitOutput(repoDir, "rev-parse", "--short", "HEAD")
	if err != nil {
		return Info{}, err
	}

	dirty := gitDirty(repoDir)

	tag, err := gitOutput(repoDir, "describe", "--tags", "--abbrev=0")
	if err != nil {
		return versionWithoutTag(repoDir, commit, dirty)
	}

	exact, _ := gitOutput(repoDir, "describe", "--tags", "--exact-match")
	commitsSince, err := gitOutput(repoDir, "rev-list", "--count", tag+"..HEAD")
	if err != nil {
		return Info{}, err
	}
	count, _ := strconv.Atoi(strings.TrimSpace(commitsSince))

	ver, err := bumpPatch(tag, count)
	if err != nil {
		return Info{}, err
	}

	return Info{
		Version:      ver,
		Commit:       shortHash(commit),
		Tag:          normalizeTag(tag),
		CommitsSince: count,
		Dirty:        dirty,
		ExactTag:     strings.TrimSpace(exact) != "",
	}, nil
}

func versionWithoutTag(repoDir, commit string, dirty bool) (Info, error) {
	total, err := gitOutput(repoDir, "rev-list", "--count", "HEAD")
	if err != nil {
		return Info{}, err
	}
	count, _ := strconv.Atoi(strings.TrimSpace(total))

	ver, err := bumpPatch("v"+DefaultBase, count)
	if err != nil {
		return Info{}, err
	}

	return Info{
		Version:      ver,
		Commit:       shortHash(commit),
		Tag:          "",
		CommitsSince: count,
		Dirty:        dirty,
		ExactTag:     false,
	}, nil
}

func bumpPatch(tag string, extra int) (string, error) {
	major, minor, patch, err := parseSemver(tag)
	if err != nil {
		return "", err
	}
	patch += extra
	return fmt.Sprintf("v%d.%d.%d", major, minor, patch), nil
}

func parseSemver(tag string) (major, minor, patch int, err error) {
	tag = strings.TrimPrefix(strings.TrimSpace(tag), "v")
	parts := strings.Split(tag, ".")
	if len(parts) < 3 {
		return 0, 0, 0, fmt.Errorf("tag semver inválida: %q", tag)
	}
	major, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, err
	}
	minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, err
	}
	patch, err = strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, err
	}
	return major, minor, patch, nil
}

func normalizeTag(tag string) string {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return ""
	}
	if strings.HasPrefix(tag, "v") {
		return tag
	}
	return "v" + tag
}

func shortHash(rev string) string {
	rev = strings.TrimSpace(rev)
	if len(rev) > 7 {
		return rev[:7]
	}
	return rev
}

func gitOutput(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func gitDirty(dir string) bool {
	cmd := exec.Command("git", "diff", "--quiet")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return true
	}
	cmd = exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = dir
	return cmd.Run() != nil
}

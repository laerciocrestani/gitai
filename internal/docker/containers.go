package docker

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// ContainerSummary describes one compose service container.
type ContainerSummary struct {
	Name    string
	Service string
	State   string
	Ports   string
	Health  string
}

// composePSRow matches docker compose ps --format json output.
type composePSRow struct {
	Name      string `json:"Name"`
	Service   string `json:"Service"`
	State     string `json:"State"`
	Status    string `json:"Status"`
	Publishers []struct {
		URL           string `json:"URL"`
		TargetPort    int    `json:"TargetPort"`
		PublishedPort int    `json:"PublishedPort"`
		Protocol      string `json:"Protocol"`
	} `json:"Publishers"`
	Health string `json:"Health"`
}

// ListComposeContainers returns containers for the compose project.
func ListComposeContainers(composeFile string) ([]ContainerSummary, error) {
	if composeFile == "" {
		return nil, nil
	}
	dir := composeDir(composeFile)
	args := []string{"compose", "-f", filepath.Base(composeFile), "ps", "--format", "json"}
	cmd := exec.Command("docker", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker compose ps: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var containers []ContainerSummary
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var row composePSRow
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			continue
		}
		state := row.State
		if state == "" {
			state = inferState(row.Status)
		}
		containers = append(containers, ContainerSummary{
			Name:    row.Name,
			Service: row.Service,
			State:   state,
			Ports:   formatPublishers(row.Publishers),
			Health:  row.Health,
		})
	}
	return containers, nil
}

func inferState(status string) string {
	lower := strings.ToLower(status)
	switch {
	case strings.Contains(lower, "up"):
		return "running"
	case strings.Contains(lower, "exit"):
		return "exited"
	case strings.Contains(lower, "pause"):
		return "paused"
	case strings.Contains(lower, "restart"):
		return "restarting"
	default:
		if status == "" {
			return "unknown"
		}
		return status
	}
}

func formatPublishers(publishers []struct {
	URL           string `json:"URL"`
	TargetPort    int    `json:"TargetPort"`
	PublishedPort int    `json:"PublishedPort"`
	Protocol      string `json:"Protocol"`
}) string {
	if len(publishers) == 0 {
		return ""
	}
	var parts []string
	for _, p := range publishers {
		if p.URL != "" {
			parts = append(parts, p.URL)
			continue
		}
		if p.PublishedPort > 0 && p.TargetPort > 0 {
			parts = append(parts, fmt.Sprintf("%d:%d", p.PublishedPort, p.TargetPort))
		}
	}
	return strings.Join(parts, ", ")
}

// ProjectName derives compose project name from directory basename.
func ProjectName(composeFile string) string {
	if composeFile == "" {
		return ""
	}
	return filepath.Base(filepath.Dir(composeFile))
}

// HasRunningContainers reports whether any container is running.
func HasRunningContainers(containers []ContainerSummary) bool {
	for _, c := range containers {
		if strings.EqualFold(c.State, "running") {
			return true
		}
	}
	return false
}

// FirstRunningService returns the first running service name.
func FirstRunningService(containers []ContainerSummary) string {
	for _, c := range containers {
		if strings.EqualFold(c.State, "running") && c.Service != "" {
			return c.Service
		}
	}
	return ""
}

package app

import (
	"fmt"
	"os"
	"strings"

	dockerpkg "github.com/laerciocrestani/openbench/internal/docker"
	"github.com/laerciocrestani/openbench/internal/ui"
)

// DockerOptions holds flags for docker commands.
type DockerOptions struct {
	ComposeFile string
	Service     string
	Build       bool
	Profile     string
	All         bool
	Tail        int
	Follow      bool
	DryRun      bool
}

func resolveComposeFile(path string) (string, error) {
	if path != "" {
		return path, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	compose := dockerpkg.FindComposeFile(cwd)
	if compose == "" {
		return "", fmt.Errorf("compose file não encontrado no diretório atual")
	}
	return compose, nil
}

// RunDockerStatus prints Docker environment status.
func RunDockerStatus() error {
	sess := ui.New("docker status", false)
	sess.Header()

	ov := dockerpkg.LoadOverview("")
	printDockerOverview(sess, ov)
	return nil
}

// RunDockerPS lists compose containers.
func RunDockerPS(opts DockerOptions) error {
	sess := ui.New("docker ps", false)
	sess.Header()

	compose, err := resolveComposeFile(opts.ComposeFile)
	if err != nil {
		return err
	}

	containers, err := dockerpkg.ListComposeContainers(compose)
	if err != nil {
		return err
	}
	if len(containers) == 0 {
		sess.Info("Nenhum container encontrado para este compose.")
		return nil
	}
	for _, c := range containers {
		line := fmt.Sprintf("%-12s %-10s %s", c.Service, c.State, c.Ports)
		if c.Health != "" {
			line += " (" + c.Health + ")"
		}
		sess.Detail(line)
	}
	return nil
}

// RunDockerUp starts compose services.
func RunDockerUp(opts DockerOptions) error {
	compose, err := resolveComposeFile(opts.ComposeFile)
	if err != nil {
		return err
	}
	if opts.DryRun {
		fmt.Printf("[dry-run] docker compose up -d (%s)\n", compose)
		return nil
	}
	return dockerpkg.Up(dockerpkg.UpOptions{
		ComposeFile: compose,
		Build:       opts.Build,
		Profile:     opts.Profile,
		DryRun:      opts.DryRun,
	})
}

// RunDockerDown stops compose services.
func RunDockerDown(opts DockerOptions) error {
	compose, err := resolveComposeFile(opts.ComposeFile)
	if err != nil {
		return err
	}
	if opts.DryRun {
		fmt.Printf("[dry-run] docker compose down (%s)\n", compose)
		return nil
	}
	return dockerpkg.Down(dockerpkg.DownOptions{
		ComposeFile: compose,
		DryRun:      opts.DryRun,
	})
}

// RunDockerLogs streams or prints service logs.
func RunDockerLogs(opts DockerOptions) error {
	compose, err := resolveComposeFile(opts.ComposeFile)
	if err != nil {
		return err
	}
	return dockerpkg.Logs(dockerpkg.LogsOptions{
		ComposeFile: compose,
		Service:     opts.Service,
		Tail:        opts.Tail,
		Follow:      opts.Follow,
	})
}

// RunDockerLogsOutput captures logs for the TUI.
func RunDockerLogsOutput(opts DockerOptions) (string, error) {
	compose, err := resolveComposeFile(opts.ComposeFile)
	if err != nil {
		return "", err
	}
	return dockerpkg.LogsOutput(dockerpkg.LogsOptions{
		ComposeFile: compose,
		Service:     opts.Service,
		Tail:        opts.Tail,
	})
}

// RunDockerShell opens an interactive shell in a service.
func RunDockerShell(opts DockerOptions) error {
	compose, err := resolveComposeFile(opts.ComposeFile)
	if err != nil {
		return err
	}
	service := opts.Service
	if service == "" {
		ov := dockerpkg.LoadOverview("")
		service = ov.DefaultService()
	}
	if service == "" {
		return fmt.Errorf("nenhum serviço em execução — informe o serviço: ob docker sh <service>")
	}
	return dockerpkg.Shell(compose, service)
}

func printDockerOverview(sess *ui.Session, ov *dockerpkg.Overview) {
	if ov == nil {
		sess.Info("Docker indisponível")
		return
	}
	sess.MetaRow("CLI", boolLabel(ov.Available))
	sess.MetaRow("Daemon", daemonLabel(ov))
	if ov.ComposeFile != "" {
		sess.MetaRow("Compose", ov.ComposeFile)
		sess.MetaRow("Project", ov.ProjectName)
	}
	if ov.Error != "" {
		sess.Warn(ov.Error)
	}
	for _, c := range ov.Containers {
		line := fmt.Sprintf("%s %s", c.Service, c.State)
		if c.Ports != "" {
			line += " " + c.Ports
		}
		if c.Health != "" {
			line += " (" + c.Health + ")"
		}
		sess.Detail(line)
	}
}

func boolLabel(ok bool) string {
	if ok {
		return "available"
	}
	return "missing"
}

func daemonLabel(ov *dockerpkg.Overview) string {
	if !ov.Available {
		return "n/a"
	}
	if ov.DaemonRunning {
		return "running"
	}
	return "stopped"
}

// CanDockerUp reports whether docker up is available from snapshot.
func CanDockerUp(snap *WorkspaceSnapshot) bool {
	return snap != nil && snap.Docker != nil && snap.Docker.CanUp() && !dockerpkg.HasRunningContainers(snap.Docker.Containers)
}

// CanDockerDown reports whether docker down is available.
func CanDockerDown(snap *WorkspaceSnapshot) bool {
	return snap != nil && snap.Docker != nil && snap.Docker.CanDown()
}

// CanDockerLogs reports whether docker logs view is available.
func CanDockerLogs(snap *WorkspaceSnapshot) bool {
	return snap != nil && snap.Docker != nil && snap.Docker.CanLogs()
}

// CanDockerShell reports whether docker shell is available.
func CanDockerShell(snap *WorkspaceSnapshot) bool {
	return snap != nil && snap.Docker != nil && snap.Docker.CanShell()
}

// DockerDefaultService returns the default service for logs/shell.
func DockerDefaultService(snap *WorkspaceSnapshot) string {
	if snap == nil || snap.Docker == nil {
		return ""
	}
	return snap.Docker.DefaultService()
}

// FormatDockerNote returns a short note for the dashboard when docker is unavailable.
func FormatDockerNote(ov *dockerpkg.Overview) string {
	if ov == nil {
		return ""
	}
	if !ov.Available {
		return "instale Docker — https://docs.docker.com/get-docker/"
	}
	if !ov.DaemonRunning {
		return "inicie o Docker daemon"
	}
	if ov.ComposeFile == "" {
		return ""
	}
	if !dockerpkg.HasRunningContainers(ov.Containers) {
		return "execute: ob docker up"
	}
	return ""
}

// DockerContainersRunning returns count of running containers.
func DockerContainersRunning(ov *dockerpkg.Overview) int {
	if ov == nil {
		return 0
	}
	n := 0
	for _, c := range ov.Containers {
		if strings.EqualFold(c.State, "running") {
			n++
		}
	}
	return n
}

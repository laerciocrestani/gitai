package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// UpOptions configures docker compose up.
type UpOptions struct {
	ComposeFile string
	Build       bool
	Profile     string
	DryRun      bool
}

// DownOptions configures docker compose down.
type DownOptions struct {
	ComposeFile string
	DryRun      bool
}

// LogsOptions configures docker compose logs.
type LogsOptions struct {
	ComposeFile string
	Service     string
	Tail        int
	Follow      bool
}

// ExecOptions configures docker compose exec.
type ExecOptions struct {
	ComposeFile string
	Service     string
	Command     []string
	Interactive bool
}

// Up runs docker compose up -d.
func Up(opts UpOptions) error {
	if opts.ComposeFile == "" {
		return fmt.Errorf("compose file não encontrado")
	}
	if opts.DryRun {
		return nil
	}
	dir := composeDir(opts.ComposeFile)
	args := []string{"compose", "-f", filepath.Base(opts.ComposeFile), "up", "-d"}
	if opts.Build {
		args = append(args, "--build")
	}
	if opts.Profile != "" {
		args = append(args, "--profile", opts.Profile)
	}
	cmd := exec.Command("docker", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Down runs docker compose down.
func Down(opts DownOptions) error {
	if opts.ComposeFile == "" {
		return fmt.Errorf("compose file não encontrado")
	}
	if opts.DryRun {
		return nil
	}
	dir := composeDir(opts.ComposeFile)
	args := []string{"compose", "-f", filepath.Base(opts.ComposeFile), "down"}
	cmd := exec.Command("docker", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Logs runs docker compose logs.
func Logs(opts LogsOptions) error {
	if opts.ComposeFile == "" {
		return fmt.Errorf("compose file não encontrado")
	}
	dir := composeDir(opts.ComposeFile)
	tail := opts.Tail
	if tail <= 0 {
		tail = 100
	}
	args := []string{"compose", "-f", filepath.Base(opts.ComposeFile), "logs", "--tail", fmt.Sprintf("%d", tail)}
	if opts.Follow {
		args = append(args, "-f")
	}
	if opts.Service != "" {
		args = append(args, opts.Service)
	}
	cmd := exec.Command("docker", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// LogsOutput captures docker compose logs without following.
func LogsOutput(opts LogsOptions) (string, error) {
	if opts.ComposeFile == "" {
		return "", fmt.Errorf("compose file não encontrado")
	}
	dir := composeDir(opts.ComposeFile)
	tail := opts.Tail
	if tail <= 0 {
		tail = 200
	}
	args := []string{"compose", "-f", filepath.Base(opts.ComposeFile), "logs", "--tail", fmt.Sprintf("%d", tail), "--no-color"}
	if opts.Service != "" {
		args = append(args, opts.Service)
	}
	cmd := exec.Command("docker", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// Exec runs a command inside a service container.
func Exec(opts ExecOptions) error {
	if opts.ComposeFile == "" {
		return fmt.Errorf("compose file não encontrado")
	}
	if opts.Service == "" {
		return fmt.Errorf("serviço não informado")
	}
	dir := composeDir(opts.ComposeFile)
	args := []string{"compose", "-f", filepath.Base(opts.ComposeFile), "exec"}
	if opts.Interactive {
		args = append(args, "-it")
	}
	args = append(args, opts.Service)
	args = append(args, opts.Command...)
	cmd := exec.Command("docker", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Shell opens an interactive shell in the service container.
func Shell(composeFile, service string) error {
	if service == "" {
		return fmt.Errorf("serviço não informado")
	}
	err := Exec(ExecOptions{
		ComposeFile: composeFile,
		Service:     service,
		Command:     []string{"sh"},
		Interactive: true,
	})
	if err == nil {
		return nil
	}
	return Exec(ExecOptions{
		ComposeFile: composeFile,
		Service:     service,
		Command:     []string{"bash"},
		Interactive: true,
	})
}

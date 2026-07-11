package docker

import (
	"os/exec"
	"runtime"
)

// HasDocker reports whether the docker CLI is available on PATH.
func HasDocker() bool {
	name := "docker"
	if runtime.GOOS == "windows" {
		name = "docker.exe"
	}
	_, err := exec.LookPath(name)
	return err == nil
}

// DaemonRunning pings the Docker daemon.
func DaemonRunning() bool {
	if !HasDocker() {
		return false
	}
	cmd := exec.Command("docker", "info")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}

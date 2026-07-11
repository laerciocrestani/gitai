package docker

import (
	"os"
	"path/filepath"
)

var composeFileNames = []string{
	"compose.yaml",
	"compose.yml",
	"docker-compose.yaml",
	"docker-compose.yml",
}

// DetectComposeFile searches for a compose file in dir.
func DetectComposeFile(dir string) string {
	dir = filepath.Clean(dir)
	for _, name := range composeFileNames {
		path := filepath.Join(dir, name)
		if fileExists(path) {
			return path
		}
	}
	return ""
}

// FindComposeFile walks from start upward looking for a compose file.
func FindComposeFile(start string) string {
	dir, err := filepath.Abs(start)
	if err != nil {
		return ""
	}
	for {
		if path := DetectComposeFile(dir); path != "" {
			return path
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func composeDir(composeFile string) string {
	return filepath.Dir(composeFile)
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

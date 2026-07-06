package git

import "testing"

func TestParsePorcelainLine(t *testing.T) {
	tests := []struct {
		line   string
		path   string
		status string
	}{
		{" M README.md", "README.md", "modified"},
		{"M  README.md", "README.md", "staged"},
		{"?? avatar.png", "avatar.png", "untracked"},
		{"MM cmd/main.go", "cmd/main.go", "staged+modified"},
	}

	for _, tc := range tests {
		path, status := parsePorcelainLine(tc.line)
		if path != tc.path || status != tc.status {
			t.Errorf("parsePorcelainLine(%q) = (%q, %q), want (%q, %q)",
				tc.line, path, status, tc.path, tc.status)
		}
	}
}

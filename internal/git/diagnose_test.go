package git

import "testing"

func TestIsBuildArtifactPath(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"node_modules/foo/bar.js", true},
		{".pnpm-store/v10/files/ab", true},
		{"sites/dist/app.js", true},
		{"sites/.astro/data-store.json", true},
		{"internal/app/doctor.go", false},
		{"src/components/Button.tsx", false},
		{"README.md", false},
	}

	for _, tt := range tests {
		if got := IsBuildArtifactPath(tt.path); got != tt.want {
			t.Errorf("IsBuildArtifactPath(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

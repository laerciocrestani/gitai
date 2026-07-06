package version

import "testing"

func TestBumpPatch(t *testing.T) {
	tests := []struct {
		tag   string
		extra int
		want  string
	}{
		{"v0.1.0", 0, "v0.1.0"},
		{"v0.1.0", 3, "v0.1.3"},
		{"v1.2.5", 1, "v1.2.6"},
	}
	for _, tc := range tests {
		got, err := bumpPatch(tc.tag, tc.extra)
		if err != nil {
			t.Fatalf("bumpPatch(%q, %d): %v", tc.tag, tc.extra, err)
		}
		if got != tc.want {
			t.Errorf("bumpPatch(%q, %d) = %q, want %q", tc.tag, tc.extra, got, tc.want)
		}
	}
}

func TestInfoDisplay(t *testing.T) {
	exact := Info{Version: "v0.1.0", ExactTag: true}
	if exact.Display() != "v0.1.0" {
		t.Errorf("exact display = %q", exact.Display())
	}

	dev := Info{Version: "v0.1.3", Commit: "3e691df", CommitsSince: 3}
	if dev.Display() != "v0.1.3 · 3e691df" {
		t.Errorf("dev display = %q", dev.Display())
	}
}

func TestParseSemver(t *testing.T) {
	major, minor, patch, err := parseSemver("v0.1.0")
	if err != nil || major != 0 || minor != 1 || patch != 0 {
		t.Fatalf("parseSemver failed: %d.%d.%d %v", major, minor, patch, err)
	}
}

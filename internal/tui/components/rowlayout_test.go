package components

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestRenderDiffBar(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		ins, del   int
		wantGreen  int
		wantRed    int
		wantGray   int
	}{
		{"empty", 0, 0, 0, 0, 5},
		{"mostly additions", 4187, 580, 4, 0, 1},
		{"balanced", 47, 48, 2, 2, 1},
		{"only additions", 10, 0, 5, 0, 0},
		{"only deletions", 0, 10, 0, 5, 0},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			out := renderDiffBar(tc.ins, tc.del)
			plain := ansi.Strip(out)
			if got := strings.Count(plain, "■"); got != tc.wantGreen+tc.wantRed {
				t.Fatalf("filled squares = %d, want %d: %q", got, tc.wantGreen+tc.wantRed, plain)
			}
			if got := strings.Count(plain, "□"); got != tc.wantGray {
				t.Fatalf("empty squares = %d, want %d: %q", got, tc.wantGray, plain)
			}
		})
	}
}

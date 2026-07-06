package ui

import "testing"

func TestIndexOf(t *testing.T) {
	opts := []string{"a", "b", "c"}
	if got := indexOf(opts, "b"); got != 1 {
		t.Fatalf("indexOf(b) = %d, want 1", got)
	}
	if got := indexOf(opts, "z"); got != -1 {
		t.Fatalf("indexOf(z) = %d, want -1", got)
	}
}

package brand

import (
	"strings"
	"testing"
)

func TestWordmarkContainsPraxicraftArt(t *testing.T) {
	w := Wordmark("Assess CLI  ·  v0.1.4")
	if !strings.Contains(w, "Assess CLI  ·  v0.1.4") {
		t.Fatalf("subtitle missing:\n%s", w)
	}
	if !strings.Contains(w, "----") {
		t.Fatal("missing frame")
	}
	// Art lines should be centered: leading spaces after | before first art glyph.
	lines := strings.Split(w, "\n")
	var artLine string
	for _, line := range lines {
		if strings.Contains(line, `____`) && strings.HasPrefix(line, "|") {
			artLine = line
			break
		}
	}
	if artLine == "" {
		t.Fatalf("art line missing:\n%s", w)
	}
	inner := strings.TrimPrefix(strings.TrimSuffix(artLine, "|"), "|")
	left := len(inner) - len(strings.TrimLeft(inner, " "))
	right := len(inner) - len(strings.TrimRight(inner, " "))
	if left < 10 {
		t.Fatalf("expected centered art with left pad, got left=%d line=%q", left, artLine)
	}
	if abs(left-right) > 2 {
		t.Fatalf("art not centered left=%d right=%d line=%q", left, right, artLine)
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func TestInteractiveIntro(t *testing.T) {
	var b strings.Builder
	InteractiveIntro(&b, "0.1.4")
	s := b.String()
	if !strings.Contains(s, "Interactive mode") {
		t.Fatal(s)
	}
	if !strings.Contains(s, CLIDocsURL) {
		t.Fatal("docs link missing")
	}
}

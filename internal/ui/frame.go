package ui

import (
	"fmt"
	"os"
	"strings"
)

const frameWidth = 72

// Rule prints a horizontal banner line (command-center style).
func Rule() {
	fmt.Fprint(os.Stdout, bannerLine())
}

// Panel prints a titled frame header.
func Panel(title string, subtitle string) {
	fmt.Fprintln(os.Stdout)
	fmt.Fprint(os.Stdout, bannerLine())
	fmt.Fprintf(os.Stdout, "  %s\n", strings.ToUpper(strings.TrimSpace(title)))
	if s := strings.TrimSpace(subtitle); s != "" {
		fmt.Fprintf(os.Stdout, "  %s\n", s)
	}
	fmt.Fprint(os.Stdout, bannerLine())
	fmt.Fprintln(os.Stdout)
}

// Footer prints the bottom hint bar.
func Footer(hints string) {
	fmt.Fprint(os.Stdout, bannerLine())
	fmt.Fprintf(os.Stdout, "  %s\n", hints)
	fmt.Fprint(os.Stdout, bannerLine())
	fmt.Fprintln(os.Stdout)
}

// OK prints a runner-style success line.
func OK(msg string) {
	fmt.Fprintf(os.Stdout, "  √  %s\n", msg)
}

// Note prints an informational line.
func Note(msg string) {
	fmt.Fprintf(os.Stdout, "  ·  %s\n", msg)
}

// Warn prints a warning line.
func Warn(msg string) {
	fmt.Fprintf(os.Stdout, "  !  %s\n", msg)
}

// Fail prints an error to stderr in frame style.
func Fail(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr)
	fmt.Fprint(os.Stderr, "  "+strings.Repeat("=", frameWidth)+"\n")
	fmt.Fprintln(os.Stderr, "  ERROR")
	fmt.Fprint(os.Stderr, "  "+strings.Repeat("=", frameWidth)+"\n")
	for _, line := range wrapWords(err.Error(), frameWidth-4) {
		fmt.Fprintf(os.Stderr, "  %s\n", line)
	}
	fmt.Fprint(os.Stderr, "  "+strings.Repeat("=", frameWidth)+"\n")
	fmt.Fprintln(os.Stderr)
}

// Done prints a soft completion strip (REPL).
func Done(msg string) {
	fmt.Fprintln(os.Stdout)
	fmt.Fprint(os.Stdout, bannerLine())
	fmt.Fprintf(os.Stdout, "  √  %s\n", msg)
	fmt.Fprint(os.Stdout, bannerLine())
	fmt.Fprintln(os.Stdout)
}

// Aborted prints a cancelled action notice.
func Aborted() {
	Note("cancelled")
}

func wrapWords(s string, width int) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if width < 20 {
		width = 20
	}
	words := strings.Fields(s)
	var lines []string
	var cur strings.Builder
	for _, w := range words {
		if cur.Len() == 0 {
			cur.WriteString(w)
			continue
		}
		if cur.Len()+1+len(w) > width {
			lines = append(lines, cur.String())
			cur.Reset()
			cur.WriteString(w)
			continue
		}
		cur.WriteByte(' ')
		cur.WriteString(w)
	}
	if cur.Len() > 0 {
		lines = append(lines, cur.String())
	}
	return lines
}

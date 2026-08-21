package brand

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const Primary = "#0D41FF"

// Banner writes the Praxicraft Assess brand header (not a MotD tips feed).
func Banner(w io.Writer, version string) {
	if w == nil {
		w = os.Stderr
	}
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(Primary)).Render("Praxicraft Assess")
	meta := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(fmt.Sprintf("CLI %s  ·  docs.praxicraft.com", version))
	fmt.Fprintf(w, "%s\n%s\n\n", title, meta)
}

// IsTTY reports whether fd is a terminal.
func IsTTY(fd uintptr) bool {
	return term.IsTerminal(int(fd))
}

// PromptPrefix is the branded REPL prompt.
func PromptPrefix() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(Primary)).Bold(true).Render("praxicraft-assess") + "> "
}

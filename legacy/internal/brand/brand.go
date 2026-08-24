package brand

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
	ProductURL     = "https://assess.praxicraft.com"
	DocsURL        = "https://docs.praxicraft.com"
	AuthDocsURL    = "https://docs.praxicraft.com/authentication"
	CLIDocsURL     = "https://docs.praxicraft.com/sdks/cli"
	APIKeysURL     = "https://assess.praxicraft.com/assess/api"
	GitHubURL      = "https://github.com/praxicraft-platform/praxicraft-assess-cli"
	DefaultBaseURL = "https://assess.praxicraft.com"
)

const wordmarkInner = 78

// VersionLine prints a simple gh-style version string.
func VersionLine(w io.Writer, version string) {
	if w == nil {
		w = os.Stdout
	}
	tag := version
	if tag != "" && !strings.HasPrefix(tag, "v") && !strings.HasPrefix(tag, "V") {
		tag = "v" + tag
	}
	fmt.Fprintf(w, "praxicraft-assess version %s\n", version)
	fmt.Fprintf(w, "%s/releases/tag/%s\n", GitHubURL, tag)
}

// Wordmark returns the line-drawn Praxicraft ASCII frame, art centered.
func Wordmark(subtitle string) string {
	bt := "`"
	art := []string{
		`____                 _                 __ _`,
		`|  _ \ _ __ __ ___  _(_) ___ _ __ __ _ / _| |_`,
		`| |_) | '__/ _` + bt + ` \ \/ / |/ __| '__/ _` + bt + ` | |_| __|`,
		`|  __/| | | (_| |>  <| | (__| | | (_| |  _| |_`,
		`|_|   |_|  \__,_/_/\_\_|\___|_|  \__,_|_|  \__|`,
	}

	sub := strings.TrimSpace(subtitle)
	if sub == "" {
		sub = "Assess CLI"
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(strings.Repeat("-", wordmarkInner+2))
	b.WriteString("\n")
	b.WriteString(frameRow(""))
	b.WriteString("\n")
	for _, line := range art {
		b.WriteString(frameRow(line))
		b.WriteString("\n")
	}
	b.WriteString(frameRow(""))
	b.WriteString("\n")
	b.WriteString(frameRow(sub))
	b.WriteString("\n")
	b.WriteString(frameRow(""))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("-", wordmarkInner+2))
	b.WriteString("\n")
	return b.String()
}

func frameRow(content string) string {
	content = strings.TrimRight(content, " ")
	n := utf8.RuneCountInString(content)
	if n > wordmarkInner {
		// Truncate by runes.
		runes := []rune(content)
		content = string(runes[:wordmarkInner])
		n = wordmarkInner
	}
	left := (wordmarkInner - n) / 2
	right := wordmarkInner - n - left
	return "|" + strings.Repeat(" ", left) + content + strings.Repeat(" ", right) + "|"
}

// ConfigureIntro prints a runner-style ASCII registration header with product links.
func ConfigureIntro(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprint(w, Wordmark("Assess CLI registration"))
	fmt.Fprint(w, `
  Praxicraft Assess is the hiring product for skills assessments, live
  interviews, pipelines, and webhooks — driven by your organisation API key.

  Product:   `+ProductURL+`
  Docs:      `+DocsURL+`
  CLI guide: `+CLIDocsURL+`
  API keys:  `+APIKeysURL+`
             (Assess → Developer → API Keys — copy ct_live_… or ct_test_… once)
  Auth help: `+AuthDocsURL+`
  Source:    `+GitHubURL+`

`)
}

// InteractiveIntro prints the wordmark for interactive shell mode.
func InteractiveIntro(w io.Writer, version string) {
	if w == nil {
		w = os.Stdout
	}
	sub := "Assess CLI"
	if version != "" {
		sub = "Assess CLI  ·  v" + strings.TrimPrefix(version, "v")
	}
	fmt.Fprint(w, Wordmark(sub))
	fmt.Fprint(w, `
  Interactive mode — press a number to open a workspace, then a letter to act.
  Docs: `+CLIDocsURL+`

`)
}

// IsTTY reports whether fd is a terminal.
func IsTTY(fd uintptr) bool {
	return term.IsTerminal(int(fd))
}

// PromptPrefix is the REPL prompt (plain text).
func PromptPrefix() string {
	return "praxicraft-assess> "
}

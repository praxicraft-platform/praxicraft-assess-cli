package brand

import (
	"fmt"
	"io"
	"os"
	"strings"

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

// ConfigureIntro prints a runner-style ASCII registration header with product links.
func ConfigureIntro(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	// Same idea as GitHub Actions ./config.sh — line-drawn wordmark, then links.
	bt := "`"
	fmt.Fprint(w, `
--------------------------------------------------------------------------------
|                                                                              |
|   ____                 _                 __ _                                |
|  |  _ \ _ __ __ ___  _(_) ___ _ __ __ _ / _| |_                              |
|  | |_) | '__/ _`+bt+` \ \/ / |/ __| '__/ _`+bt+` | |_| __|                             |
|  |  __/| | | (_| |>  <| | (__| | | (_| |  _| |_                              |
|  |_|   |_|  \__,_/_/\_\_|\___|_|  \__,_|_|  \__|                             |
|                                                                              |
|                     Assess CLI registration                                  |
|                                                                              |
--------------------------------------------------------------------------------

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

// IsTTY reports whether fd is a terminal.
func IsTTY(fd uintptr) bool {
	return term.IsTerminal(int(fd))
}

// PromptPrefix is the REPL prompt (plain text).
func PromptPrefix() string {
	return "praxicraft-assess> "
}

package ui

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/brand"
)

// SessionInfo is optional context shown on the interactive home screen.
type SessionInfo struct {
	Version string
	Profile string
	BaseURL string
	OrgName string // optional; empty if unknown
}

// Shell opens the nested interactive experience and returns a CLI command line to run.
func Shell(opts Options, session SessionInfo, roots []MenuNode) (string, error) {
	if !opts.Interactive {
		return "", &api.UsageError{Msg: "interactive shell requires a TTY"}
	}
	if !canUseTUI(opts) {
		return NavigateMenu(opts, "Choose a resource", roots)
	}

	m := shellModel{
		session: session,
		stack:   []shellFrame{{title: "home", nodes: roots}},
	}
	p := tea.NewProgram(m, tea.WithOutput(os.Stdout), tea.WithInput(os.Stdin))
	final, err := p.Run()
	if err != nil {
		return "", err
	}
	out, ok := final.(shellModel)
	if !ok {
		return "", &api.UsageError{Msg: "shell failed"}
	}
	if out.abort {
		return "", &api.AbortError{}
	}
	if out.command == "" {
		return "", &api.AbortError{}
	}
	return out.command, nil
}

type shellFrame struct {
	title string
	nodes []MenuNode
}

type shellModel struct {
	session SessionInfo
	stack   []shellFrame
	cursor  int
	command string
	abort   bool
	quit    bool
}

func (m shellModel) Init() tea.Cmd { return nil }

func (m shellModel) cur() shellFrame {
	return m.stack[len(m.stack)-1]
}

func (m shellModel) atHome() bool {
	return len(m.stack) == 1
}

func (m shellModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	nodes := m.cur().nodes
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "q":
			if m.atHome() {
				m.command = "exit"
				m.quit = true
				return m, tea.Quit
			}
			m.abort = true
			m.quit = true
			return m, tea.Quit
		case "esc", "left", "backspace":
			if !m.atHome() {
				m.stack = m.stack[:len(m.stack)-1]
				m.cursor = 0
				return m, nil
			}
			m.command = "exit"
			m.quit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(nodes)-1 {
				m.cursor++
			}
		case "enter", "right", "l":
			return m.activate(m.cursor)
		default:
			// Number or letter shortcut
			if n, ok := shortcutIndex(key, nodes); ok {
				return m.activate(n)
			}
		}
	}
	return m, nil
}

func (m shellModel) activate(i int) (tea.Model, tea.Cmd) {
	nodes := m.cur().nodes
	if i < 0 || i >= len(nodes) {
		return m, nil
	}
	n := nodes[i]
	if len(n.Children) > 0 {
		m.stack = append(m.stack, shellFrame{title: n.Label, nodes: n.Children})
		m.cursor = 0
		return m, nil
	}
	cmd := strings.TrimSpace(n.Command)
	if cmd == "" {
		return m, nil
	}
	m.command = cmd
	m.quit = true
	return m, tea.Quit
}

func shortcutIndex(key string, nodes []MenuNode) (int, bool) {
	if key == "" {
		return 0, false
	}
	// Digits 1-9 (and 0 as 10th)
	if len(key) == 1 && key[0] >= '1' && key[0] <= '9' {
		i := int(key[0] - '1')
		if i < len(nodes) {
			return i, true
		}
	}
	if key == "0" && len(nodes) >= 10 {
		return 9, true
	}
	// Letter keys match MenuNode.Key (case-insensitive)
	k := strings.ToLower(key)
	if len(k) != 1 || !unicode.IsLetter(rune(k[0])) {
		return 0, false
	}
	for i, n := range nodes {
		if strings.ToLower(n.Key) == k {
			return i, true
		}
	}
	return 0, false
}

func (m shellModel) View() string {
	var b strings.Builder
	if m.atHome() {
		m.renderHome(&b)
	} else {
		m.renderResource(&b)
	}
	return b.String()
}

func (m shellModel) renderHome(b *strings.Builder) {
	b.WriteString("\n")
	b.WriteString(bannerLine())
	b.WriteString("  COMMAND CENTER")
	if ver := strings.TrimPrefix(m.session.Version, "v"); ver != "" {
		fmt.Fprintf(b, "  ·  v%s", ver)
	}
	b.WriteString("\n")
	b.WriteString(bannerLine())
	b.WriteString("\n")

	prof := m.session.Profile
	if prof == "" {
		prof = "default"
	}
	base := m.session.BaseURL
	if base == "" {
		base = brand.DefaultBaseURL
	}
	base = strings.TrimPrefix(strings.TrimPrefix(base, "https://"), "http://")
	if m.session.OrgName != "" {
		fmt.Fprintf(b, "  org      %s\n", m.session.OrgName)
	}
	fmt.Fprintf(b, "  session  profile %s  ·  %s\n", prof, base)
	b.WriteString("\n")
	b.WriteString("  Pick a workspace (press the number — instant):\n\n")

	nodes := m.cur().nodes
	for i, n := range nodes {
		mark := " "
		if i == m.cursor {
			mark = ">"
		}
		key := n.Key
		if key == "" {
			key = fmt.Sprintf("%d", i+1)
		}
		fmt.Fprintf(b, "  %s [%s]  %-14s    %s\n", mark, key, n.Label, n.Hint)
		if i == m.cursor {
			b.WriteString("      " + strings.Repeat("·", 52) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(bannerLine())
	b.WriteString("  1–9 open workspace  ·  ↑/↓ move  ·  Enter  ·  q quit\n")
	b.WriteString(bannerLine())
	b.WriteString("\n")
}

func (m shellModel) renderResource(b *strings.Builder) {
	cur := m.cur()
	b.WriteString("\n")
	b.WriteString(bannerLine())
	fmt.Fprintf(b, "  %s\n", strings.ToUpper(cur.title))
	path := "home"
	for _, f := range m.stack[1:] {
		path += " › " + f.title
	}
	fmt.Fprintf(b, "  %s\n", path)
	b.WriteString(bannerLine())
	b.WriteString("\n")
	b.WriteString("  Pick an action (press the letter — instant):\n\n")

	for i, n := range cur.nodes {
		mark := " "
		if i == m.cursor {
			mark = ">"
		}
		key := n.Key
		if key == "" {
			key = string(rune('a' + i))
		}
		fmt.Fprintf(b, "  %s [%s]  %-14s    %s\n", mark, strings.ToUpper(key), n.Label, n.Hint)
		if i == m.cursor {
			b.WriteString("      " + strings.Repeat("·", 52) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(bannerLine())
	b.WriteString("  letter runs  ·  ← / esc back to workspaces  ·  ↑/↓  ·  q quit\n")
	b.WriteString(bannerLine())
	b.WriteString("\n")
}

func bannerLine() string {
	return "  " + strings.Repeat("=", 72) + "\n"
}

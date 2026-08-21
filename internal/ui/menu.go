package ui

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/brand"
)

// Choice is a selectable option (arrow keys + Enter).
type Choice struct {
	Label string
	Value string
	Hint  string
	Key   string // optional instant key (y/n, digits)
}

const (
	viewWindow = 14
	cmdCol     = 24
)

type menuModel struct {
	title     string
	subtitle  string
	choices   []Choice
	cursor    int
	offset    int
	picked    string
	quit      bool
	abort     bool
	backValue string
}

func (m menuModel) Init() tea.Cmd { return nil }

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "q":
			m.abort = true
			m.quit = true
			return m, tea.Quit
		case "esc", "left", "backspace":
			if m.backValue != "" {
				m.picked = m.backValue
				m.quit = true
				return m, tea.Quit
			}
			m.abort = true
			m.quit = true
			return m, tea.Quit
		case "enter", "right":
			if len(m.choices) == 0 {
				m.abort = true
				m.quit = true
				return m, tea.Quit
			}
			m.picked = m.choices[m.cursor].Value
			m.quit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.offset {
					m.offset = m.cursor
				}
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
				if m.cursor >= m.offset+viewWindow {
					m.offset = m.cursor - viewWindow + 1
				}
			}
		case "pgup":
			m.cursor = max(0, m.cursor-viewWindow)
			m.offset = max(0, m.offset-viewWindow)
		case "pgdown":
			m.cursor = min(len(m.choices)-1, m.cursor+viewWindow)
			if m.cursor >= m.offset+viewWindow {
				m.offset = m.cursor - viewWindow + 1
			}
		case "home", "g":
			m.cursor = 0
			m.offset = 0
		case "end", "G":
			if len(m.choices) > 0 {
				m.cursor = len(m.choices) - 1
				m.offset = max(0, len(m.choices)-viewWindow)
			}
		default:
			if i, ok := choiceKeyIndex(key, m.choices); ok {
				m.picked = m.choices[i].Value
				m.quit = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func choiceKeyIndex(key string, choices []Choice) (int, bool) {
	if len(key) != 1 {
		return 0, false
	}
	k := strings.ToLower(key)
	for i, c := range choices {
		if c.Key != "" && strings.ToLower(c.Key) == k {
			return i, true
		}
	}
	if key[0] >= '1' && key[0] <= '9' {
		i := int(key[0] - '1')
		if i < len(choices) {
			return i, true
		}
	}
	return 0, false
}

func (m menuModel) View() string {
	var b strings.Builder
	title := strings.TrimSpace(m.title)
	if title == "" {
		title = "Select"
	}
	b.WriteString("\n")
	b.WriteString(bannerLine())
	fmt.Fprintf(&b, "  %s\n", strings.ToUpper(title))
	if sub := strings.TrimSpace(m.subtitle); sub != "" {
		fmt.Fprintf(&b, "  %s\n", sub)
	}
	b.WriteString(bannerLine())
	b.WriteString("\n")

	if len(m.choices) == 0 {
		b.WriteString("  ·  nothing to show\n\n")
	}

	end := m.offset + viewWindow
	if end > len(m.choices) {
		end = len(m.choices)
	}
	for i := m.offset; i < end; i++ {
		c := m.choices[i]
		line := formatChoiceLine(c)
		mark := " "
		if i == m.cursor {
			mark = ">"
		}
		key := c.Key
		if key == "" && i < 9 {
			key = fmt.Sprintf("%d", i+1)
		}
		if key != "" {
			fmt.Fprintf(&b, "  %s [%s]  %s\n", mark, strings.ToUpper(key), line)
		} else {
			fmt.Fprintf(&b, "  %s      %s\n", mark, line)
		}
		if i == m.cursor {
			b.WriteString("      " + strings.Repeat("·", 52) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(bannerLine())
	hints := "↑/↓ move · Enter select · esc cancel"
	if m.backValue != "" {
		hints = "↑/↓ · Enter · ← / esc back · q quit"
	}
	if m.offset > 0 || end < len(m.choices) {
		hints += fmt.Sprintf(" · %d–%d/%d", m.offset+1, end, len(m.choices))
	}
	fmt.Fprintf(&b, "  %s\n", hints)
	b.WriteString(bannerLine())
	b.WriteString("\n")
	return b.String()
}

func formatChoiceLine(c Choice) string {
	cmd := strings.TrimSpace(c.Label)
	if cmd == "" {
		cmd = c.Value
	}
	hint := strings.TrimSpace(c.Hint)
	if hint == "" {
		if left, right, ok := strings.Cut(cmd, " — "); ok {
			cmd, hint = left, right
		}
	}
	if hint == "" {
		return cmd
	}
	return padRight(cmd, cmdCol) + "  " + hint
}

func padRight(s string, width int) string {
	n := utf8.RuneCountInString(s)
	if n >= width {
		return s
	}
	return s + strings.Repeat(" ", width-n)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// SelectRunner shows a command-center style ↑/↓ menu.
func SelectRunner(opts Options, title string, choices []Choice) (string, error) {
	return runSelect(opts, title, "", choices, "")
}

// SelectWithHint is Select with a subtitle under the panel title.
func SelectWithHint(opts Options, title, subtitle string, choices []Choice) (string, error) {
	return runSelect(opts, title, subtitle, choices, "")
}

func selectRunner(opts Options, title string, choices []Choice, backValue string) (string, error) {
	return runSelect(opts, title, "", choices, backValue)
}

func runSelect(opts Options, title, subtitle string, choices []Choice, backValue string) (string, error) {
	if len(choices) == 0 {
		return "", &api.UsageError{Msg: "nothing to select"}
	}
	if len(choices) == 1 && !strings.HasPrefix(choices[0].Value, "__") {
		return choices[0].Value, nil
	}
	if !opts.Interactive {
		return "", &api.UsageError{Msg: title + " requires an argument (or run interactively)"}
	}

	if canUseTUI(opts) {
		m := menuModel{title: title, subtitle: subtitle, choices: choices, backValue: backValue}
		p := tea.NewProgram(m, tea.WithOutput(os.Stdout), tea.WithInput(os.Stdin))
		final, err := p.Run()
		if err != nil {
			return "", err
		}
		out, ok := final.(menuModel)
		if !ok {
			return "", &api.UsageError{Msg: "menu failed"}
		}
		if out.abort || out.picked == "" {
			return "", &api.AbortError{}
		}
		return out.picked, nil
	}

	Panel(title, subtitle)
	for i, c := range choices {
		fmt.Fprintf(os.Stdout, "  %2d) %s\n", i+1, formatChoiceLine(c))
	}
	fmt.Fprintln(os.Stdout)
	line, err := PromptEnter(opts, "Enter number:", "")
	if err != nil {
		return "", err
	}
	var n int
	if _, err := fmt.Sscanf(strings.TrimSpace(line), "%d", &n); err != nil || n < 1 || n > len(choices) {
		return "", &api.UsageError{Msg: "invalid selection"}
	}
	return choices[n-1].Value, nil
}

// Select uses the framed menu when a TTY is available.
func Select(opts Options, title string, choices []Choice) (string, error) {
	return SelectRunner(opts, title, choices)
}

func canUseTUI(opts Options) bool {
	return opts.Interactive && brand.IsTTY(os.Stdin.Fd()) && brand.IsTTY(os.Stdout.Fd())
}

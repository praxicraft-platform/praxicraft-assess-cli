package ui

import (
	"fmt"
	"strings"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
)

const (
	navBack = "__back__"
	navQuit = "__quit__"
)

// NavigateMenu is a non-TUI / fallback nested picker (numbers only).
func NavigateMenu(opts Options, rootTitle string, roots []MenuNode) (string, error) {
	if !opts.Interactive {
		return "", &api.UsageError{Msg: "interactive menu requires a TTY (omit --non-interactive)"}
	}

	type frame struct {
		title string
		nodes []MenuNode
	}
	stack := []frame{{title: rootTitle, nodes: roots}}

	for {
		cur := stack[len(stack)-1]
		choices := make([]Choice, 0, len(cur.nodes)+2)
		for _, n := range cur.nodes {
			label := n.Label
			if n.Key != "" {
				label = fmt.Sprintf("[%s] %s", n.Key, n.Label)
			}
			choices = append(choices, Choice{Label: label, Hint: n.Hint, Value: n.Label})
		}
		back := ""
		if len(stack) > 1 {
			choices = append(choices, Choice{Label: "← Back", Hint: "previous menu", Value: navBack})
			back = navBack
		}
		choices = append(choices, Choice{Label: "Exit", Hint: "leave interactive mode", Value: navQuit})

		picked, err := selectRunner(opts, cur.title, choices, back)
		if err != nil {
			return "", err
		}
		switch picked {
		case navQuit:
			return "exit", nil
		case navBack:
			if len(stack) == 1 {
				continue
			}
			stack = stack[:len(stack)-1]
			continue
		}

		var node *MenuNode
		for i := range cur.nodes {
			if cur.nodes[i].Label == picked {
				node = &cur.nodes[i]
				break
			}
		}
		if node == nil {
			return "", &api.UsageError{Msg: "unknown menu item"}
		}
		if len(node.Children) > 0 {
			stack = append(stack, frame{title: node.Label, nodes: node.Children})
			continue
		}
		if strings.TrimSpace(node.Command) == "" {
			return "", &api.UsageError{Msg: fmt.Sprintf("%s has no action", node.Label)}
		}
		return node.Command, nil
	}
}

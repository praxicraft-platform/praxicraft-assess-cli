package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/brand"
)

// Choice is a selectable option (arrow keys + Enter).
type Choice struct {
	Label string
	Value string
}

func canUseTUI(opts Options) bool {
	return opts.Interactive && brand.IsTTY(os.Stdin.Fd()) && brand.IsTTY(os.Stdout.Fd())
}

// Select shows an arrow-key menu. Falls back to numbered list when not a TTY.
// A single real choice is returned immediately unless it is a sentinel (e.g. load more).
func Select(opts Options, title string, choices []Choice) (string, error) {
	if len(choices) == 0 {
		return "", &api.UsageError{Msg: "nothing to select"}
	}
	if len(choices) == 1 && !strings.HasPrefix(choices[0].Value, "__praxicraft_") {
		return choices[0].Value, nil
	}
	if !opts.Interactive {
		return "", &api.UsageError{Msg: title + " requires an argument (or run interactively)"}
	}

	if canUseTUI(opts) {
		options := make([]huh.Option[string], 0, len(choices))
		for _, c := range choices {
			label := c.Label
			if label == "" {
				label = c.Value
			}
			options = append(options, huh.NewOption(label, c.Value))
		}
		var selected string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title(title).
					Description("↑/↓ move · Enter select · Ctrl+C cancel").
					Options(options...).
					Value(&selected),
			),
		).WithTheme(huh.ThemeBase())
		if err := form.Run(); err != nil {
			return "", err
		}
		if selected == "" {
			return "", &api.AbortError{}
		}
		return selected, nil
	}

	// Non-TUI fallback: numbered list.
	fmt.Fprintln(os.Stdout, title)
	for i, c := range choices {
		label := c.Label
		if label == "" {
			label = c.Value
		}
		fmt.Fprintf(os.Stdout, "  %d) %s\n", i+1, label)
	}
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

// ConfirmTUI asks Yes/No with arrow keys when possible.
func ConfirmTUI(opts Options, message string) (bool, error) {
	if opts.Yes {
		return true, nil
	}
	if !opts.Interactive {
		return false, &api.UsageError{Msg: "refusing destructive action without --yes in non-interactive mode"}
	}
	if canUseTUI(opts) {
		var ok bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(message).
					Affirmative("Yes").
					Negative("No").
					Value(&ok),
			),
		).WithTheme(huh.ThemeBase())
		if err := form.Run(); err != nil {
			return false, err
		}
		return ok, nil
	}
	return Confirm(opts, message)
}

// InputTUI prompts for a single line with Enter to submit.
func InputTUI(opts Options, title, placeholder, def string) (string, error) {
	if !opts.Interactive {
		if strings.TrimSpace(def) != "" {
			return strings.TrimSpace(def), nil
		}
		return "", requiredErr(title)
	}
	if canUseTUI(opts) {
		value := def
		input := huh.NewInput().Title(title).Value(&value)
		if placeholder != "" {
			input = input.Placeholder(placeholder)
		}
		form := huh.NewForm(huh.NewGroup(input)).WithTheme(huh.ThemeBase())
		if err := form.Run(); err != nil {
			return "", err
		}
		value = strings.TrimSpace(value)
		if value == "" {
			return "", requiredErr(title)
		}
		return value, nil
	}
	return PromptEnter(opts, title+":", def)
}

// FormInvite walks invite create fields with arrow/enter UX.
func FormInvite(opts Options, email, name string, sendEmail bool) (string, string, bool, error) {
	if !opts.Interactive {
		if email == "" {
			return "", "", false, requiredErr("Candidate email")
		}
		return email, name, sendEmail, nil
	}
	if canUseTUI(opts) {
		e, n, send := email, name, sendEmail
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Candidate email").Placeholder("candidate@example.com").Value(&e),
				huh.NewInput().Title("Candidate name (optional)").Value(&n),
				huh.NewConfirm().Title("Send invite email?").Affirmative("Yes").Negative("No").Value(&send),
			),
		).WithTheme(huh.ThemeBase())
		if err := form.Run(); err != nil {
			return "", "", false, err
		}
		e = strings.TrimSpace(e)
		if e == "" {
			return "", "", false, requiredErr("Candidate email")
		}
		return e, strings.TrimSpace(n), send, nil
	}
	e, err := PromptString(opts, "Candidate email", email)
	if err != nil {
		return "", "", false, err
	}
	n := name
	if n == "" {
		n, _ = PromptEnter(opts, "Enter Candidate name (optional):", "")
	}
	send := sendEmail
	ok, err := Confirm(opts, "Send invite email?")
	if err != nil {
		return "", "", false, err
	}
	send = ok
	return e, strings.TrimSpace(n), send, nil
}

package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
)

// ConfirmTUI asks Yes/No with a framed panel (Y/N instant keys).
func ConfirmTUI(opts Options, message string) (bool, error) {
	if opts.Yes {
		return true, nil
	}
	if !opts.Interactive {
		return false, &api.UsageError{Msg: "refusing destructive action without --yes in non-interactive mode"}
	}
	if canUseTUI(opts) {
		v, err := runSelect(opts, "Confirm", message, []Choice{
			{Key: "y", Label: "Yes", Hint: "continue", Value: "yes"},
			{Key: "n", Label: "No", Hint: "cancel", Value: "no"},
		}, "")
		if err != nil {
			return false, err
		}
		return v == "yes", nil
	}
	return Confirm(opts, message)
}

// InputTUI prompts for a single line (runner-style Enter).
func InputTUI(opts Options, title, placeholder, def string) (string, error) {
	_ = placeholder
	return PromptEnter(opts, title+":", def)
}

// FormInvite walks invite create fields in framed wizard steps.
func FormInvite(opts Options, email, name string, sendEmail bool) (string, string, bool, error) {
	if !opts.Interactive {
		if email == "" {
			return "", "", false, requiredErr("Candidate email")
		}
		return email, name, sendEmail, nil
	}

	Panel("Create invite", "step 1 of 3 · candidate email")
	e, err := PromptString(opts, "Candidate email", email)
	if err != nil {
		return "", "", false, err
	}

	Panel("Create invite", "step 2 of 3 · name (optional)")
	n := name
	if n == "" {
		fmt.Fprint(os.Stdout, "Enter Candidate name (optional): ")
		line, _ := readLine()
		n = strings.TrimSpace(line)
	} else {
		Note("using --name " + n)
	}

	Panel("Create invite", "step 3 of 3 · send email?")
	send := sendEmail
	if canUseTUI(opts) {
		v, err := runSelect(opts, "Send invite email?", "notify the candidate now", []Choice{
			{Key: "y", Label: "Yes", Hint: "send email", Value: "yes"},
			{Key: "n", Label: "No", Hint: "create invite only", Value: "no"},
		}, "")
		if err != nil {
			return "", "", false, err
		}
		send = v == "yes"
	} else {
		ok, err := Confirm(opts, "Send invite email?")
		if err != nil {
			return "", "", false, err
		}
		send = ok
	}
	return e, strings.TrimSpace(n), send, nil
}

package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
)

// Enabled reports whether interactive prompts are allowed.
type Options struct {
	Interactive bool
	Yes         bool
}

// Confirm asks y/n unless Yes or non-interactive.
func Confirm(opts Options, message string) (bool, error) {
	if opts.Yes {
		return true, nil
	}
	if !opts.Interactive {
		return false, fmt.Errorf("refusing destructive action without --yes in non-interactive mode")
	}
	var ok bool
	err := huh.NewConfirm().Title(message).Value(&ok).Run()
	return ok, err
}

// PromptString asks for a required string when empty and interactive.
func PromptString(opts Options, title, current string) (string, error) {
	if strings.TrimSpace(current) != "" {
		return current, nil
	}
	if !opts.Interactive {
		return "", fmt.Errorf("%s is required", title)
	}
	var v string
	err := huh.NewInput().Title(title).Value(&v).Run()
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(v) == "" {
		return "", fmt.Errorf("%s is required", title)
	}
	return strings.TrimSpace(v), nil
}

// PromptSecret asks for an API key.
func PromptSecret(opts Options, title string) (string, error) {
	if !opts.Interactive {
		return "", fmt.Errorf("%s is required", title)
	}
	var v string
	err := huh.NewInput().Title(title).EchoMode(huh.EchoModePassword).Value(&v).Run()
	return strings.TrimSpace(v), err
}

// ReadLine reads a single line from stdin (REPL).
func ReadLine(prompt string) (string, error) {
	fmt.Fprint(os.Stdout, prompt)
	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		if err := sc.Err(); err != nil {
			return "", err
		}
		return "", ioEOF
	}
	return sc.Text(), nil
}

var ioEOF = fmt.Errorf("EOF")

// IsEOF reports end of REPL input.
func IsEOF(err error) bool {
	return err == ioEOF
}

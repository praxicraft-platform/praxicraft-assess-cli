package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
	"golang.org/x/term"
)

// Options controls interactive prompting.
type Options struct {
	Interactive bool
	Yes         bool
}

var (
	stdinOnce sync.Once
	stdinR    *bufio.Reader
)

func stdin() *bufio.Reader {
	stdinOnce.Do(func() {
		stdinR = bufio.NewReader(os.Stdin)
	})
	return stdinR
}

// Confirm asks y/N unless Yes or non-interactive (runner-style).
func Confirm(opts Options, message string) (bool, error) {
	if opts.Yes {
		return true, nil
	}
	if !opts.Interactive {
		return false, &api.UsageError{Msg: "refusing destructive action without --yes in non-interactive mode"}
	}
	fmt.Fprintf(os.Stdout, "%s (Y/N) [press Enter for N] ", message)
	line, err := readLine()
	if err != nil {
		return false, err
	}
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}

// PromptEnter asks a question in GitHub Actions runner style.
// Example: Enter the name of the runner: [press Enter for hostname]
func PromptEnter(opts Options, question, def string) (string, error) {
	if !opts.Interactive {
		if strings.TrimSpace(def) != "" {
			return strings.TrimSpace(def), nil
		}
		return "", requiredErr(question)
	}
	if def != "" {
		fmt.Fprintf(os.Stdout, "%s [press Enter for %s] ", question, def)
	} else {
		fmt.Fprintf(os.Stdout, "%s ", question)
	}
	line, err := readLine()
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		line = strings.TrimSpace(def)
	}
	if line == "" {
		return "", requiredErr(question)
	}
	return line, nil
}

// PromptSecretEnter asks for a secret with no echo (runner-style question text).
func PromptSecretEnter(opts Options, question string) (string, error) {
	if !opts.Interactive {
		return "", requiredErr(question)
	}
	fmt.Fprintf(os.Stdout, "%s ", question)
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		b, err := term.ReadPassword(fd)
		fmt.Fprintln(os.Stdout)
		if err != nil {
			return "", err
		}
		v := strings.TrimSpace(string(b))
		if v == "" {
			return "", requiredErr(question)
		}
		return v, nil
	}
	line, err := readLine()
	if err != nil {
		return "", err
	}
	v := strings.TrimSpace(line)
	if v == "" {
		return "", requiredErr(question)
	}
	return v, nil
}

func requiredErr(question string) error {
	q := strings.TrimSpace(question)
	q = strings.TrimSuffix(q, ":")
	return &api.UsageError{Msg: q + " is required (pass a flag, or omit --non-interactive)"}
}

// OK prints a runner-style success line.
func OK(msg string) {
	fmt.Fprintf(os.Stdout, "√ %s\n", msg)
}

// Section prints a runner-style section header.
func Section(title string) {
	fmt.Fprintf(os.Stdout, "\n# %s\n\n", title)
}

// PromptString keeps a short label form for missing command flags.
func PromptString(opts Options, label, def string) (string, error) {
	q := "Enter " + label + ":"
	return PromptEnter(opts, q, def)
}

// PromptSecret keeps a short label form for missing command flags.
func PromptSecret(opts Options, label string) (string, error) {
	return PromptSecretEnter(opts, "Enter "+label+":")
}

// ReadLine reads a single line from stdin (REPL).
func ReadLine(prompt string) (string, error) {
	fmt.Fprint(os.Stdout, prompt)
	return readLine()
}

func readLine() (string, error) {
	line, err := stdin().ReadString('\n')
	if err != nil {
		if err == io.EOF {
			if strings.TrimSpace(line) != "" {
				return strings.TrimRight(line, "\r\n"), nil
			}
			return "", ioEOF
		}
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

var ioEOF = fmt.Errorf("EOF")

// IsEOF reports end of REPL input.
func IsEOF(err error) bool {
	return err == ioEOF
}

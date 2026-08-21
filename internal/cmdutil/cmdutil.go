package cmdutil

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/ui"
	"github.com/spf13/cobra"
)

// Run wraps an API call with EnsureAPI + Print + exit-friendly errors.
func Run(rt *runtime.Runtime, fn func() (any, error)) error {
	if err := rt.EnsureAPI(); err != nil {
		return err
	}
	out, err := fn()
	if err != nil {
		return err
	}
	if out == nil {
		return nil
	}
	return rt.Print(out)
}

// BodyFlag adds --body for JSON payloads.
func BodyFlag(cmd *cobra.Command, dest *string) {
	cmd.Flags().StringVar(dest, "body", "", "JSON request body")
}

// ReadBodyFile loads --body or --body-file.
func ReadBody(body, bodyFile string) (map[string]any, error) {
	raw := body
	if bodyFile != "" {
		b, err := os.ReadFile(bodyFile)
		if err != nil {
			return nil, err
		}
		raw = string(b)
	}
	return runtime.ParseBody(raw)
}

// BodyFlags adds --body and --body-file.
func BodyFlags(cmd *cobra.Command, body, bodyFile *string) {
	cmd.Flags().StringVar(body, "body", "", "JSON request body")
	cmd.Flags().StringVar(bodyFile, "body-file", "", "Path to JSON request body file")
}

// QueryFromPairs builds url.Values from key=value pairs.
func QueryFromPairs(pairs []string) url.Values {
	q := url.Values{}
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			continue
		}
		q.Add(k, v)
	}
	return q
}

// ConfirmDestructive confirms unless --yes.
func ConfirmDestructive(rt *runtime.Runtime, msg string) error {
	ok, err := ui.Confirm(rt.UI, msg)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("aborted")
	}
	return nil
}

// ExitError wraps API errors for main.
func ExitError(err error) int {
	return api.ExitCode(err)
}

// MustJSON pretty helper for tests.
func MustJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

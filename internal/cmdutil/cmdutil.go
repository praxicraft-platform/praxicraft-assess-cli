package cmdutil

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
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

// ReadBody loads --body or --body-file.
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

// FilterFlag adds --filter key=value for API list query params (not JMESPath).
func FilterFlag(cmd *cobra.Command, dest *[]string) {
	cmd.Flags().StringArrayVar(dest, "filter", nil, "API query filter as key=value (repeatable)")
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

// ExitError wraps API errors for main.
func ExitError(err error) int {
	return api.ExitCode(err)
}

// MustJSON pretty helper for tests.
func MustJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

// RequireBody returns a usage error when body is empty and interactive create needs JSON.
func RequireBody(body map[string]any, hint string) error {
	if len(body) == 0 {
		if hint == "" {
			hint = "pass --body '{...}' or --body-file path.json"
		}
		return &api.UsageError{Msg: hint}
	}
	return nil
}

// ErrfUsage returns a usage error.
func ErrfUsage(format string, args ...any) error {
	return &api.UsageError{Msg: fmt.Sprintf(format, args...)}
}

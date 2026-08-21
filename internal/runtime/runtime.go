package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/brand"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/config"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/output"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/ui"
)

// Options are global CLI flags.
type Options struct {
	Profile        string
	APIKey         string
	BaseURL        string
	Output         string
	Query          string // JMESPath (--query)
	Yes            bool
	NonInteractive bool
	NoBanner       bool
	Version        string
}

// Runtime is the shared session for commands.
type Runtime struct {
	Opts Options
	Out  *output.Printer
	UI   ui.Options
	API  *api.Client
}

// New builds a Runtime and resolves config. API client is created lazily via EnsureAPI.
func New(opts Options) *Runtime {
	// Prompts run whenever --non-interactive is not set (same idea as runner config.sh).
	// Output format still prefers table only on a real TTY.
	promptsOK := !opts.NonInteractive
	format := opts.Output
	if format == "" {
		if promptsOK && brand.IsTTY(os.Stdout.Fd()) {
			format = string(output.Table)
		} else {
			format = string(output.JSON)
		}
	}
	return &Runtime{
		Opts: opts,
		Out:  output.NewPrinter(format, opts.Query),
		UI:   ui.Options{Interactive: promptsOK, Yes: opts.Yes},
	}
}

// EnsureAPI resolves credentials and constructs the API client.
func (r *Runtime) EnsureAPI() error {
	if r.API != nil {
		return nil
	}
	res, err := config.Resolve(r.Opts.Profile, r.Opts.APIKey, r.Opts.BaseURL)
	if err != nil {
		return &api.UsageError{Msg: err.Error()}
	}
	c, err := api.New(res.APIKey, res.BaseURL)
	if err != nil {
		return err
	}
	if r.Opts.Version != "" {
		c.UserAgent = "praxicraft-assess-cli/" + r.Opts.Version
	}
	r.API = c
	return nil
}

// Print writes a value using the configured printer.
func (r *Runtime) Print(v any) error {
	return r.Out.Print(v)
}

// ParseBody unmarshals --body JSON into map.
func ParseBody(raw string) (map[string]any, error) {
	if raw == "" {
		return map[string]any{}, nil
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, err
	}
	return m, nil
}

// Context returns a background context (hook for future deadlines).
func (r *Runtime) Context() context.Context {
	return context.Background()
}

// ShowBanner prints a quiet interactive welcome (gh-style), not a product MotD.
func (r *Runtime) ShowBanner() {
	if r.Opts.NoBanner || r.Opts.NonInteractive {
		return
	}
	if !brand.IsTTY(os.Stderr.Fd()) {
		return
	}
	fmt.Fprintf(os.Stderr, "praxicraft-assess %s\nType 'help' for commands, or 'exit' to quit.\n\n", r.Opts.Version)
}

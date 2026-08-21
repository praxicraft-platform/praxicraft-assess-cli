package runtime

import (
	"context"
	"encoding/json"
	"os"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/brand"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/config"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/output"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/ui"
)

// Options are global CLI flags.
type Options struct {
	Profile         string
	APIKey          string
	BaseURL         string
	Output          string
	Query           string
	Yes             bool
	NonInteractive  bool
	NoBanner        bool
	Debug           bool
	MaxItems        int
	StartingToken   string
	NoPaginate      bool
	Version         string
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
	interactive := !opts.NonInteractive && brand.IsTTY(os.Stdin.Fd()) && brand.IsTTY(os.Stdout.Fd())
	format := opts.Output
	if format == "" {
		if interactive {
			format = string(output.Table)
		} else {
			format = string(output.JSON)
		}
	}
	return &Runtime{
		Opts: opts,
		Out:  output.NewPrinter(format, opts.Query),
		UI:   ui.Options{Interactive: interactive, Yes: opts.Yes},
	}
}

// EnsureAPI resolves credentials and constructs the API client.
func (r *Runtime) EnsureAPI() error {
	if r.API != nil {
		return nil
	}
	res, err := config.Resolve(r.Opts.Profile, r.Opts.APIKey, r.Opts.BaseURL)
	if err != nil {
		return err
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

// ShowBanner prints brand banner when appropriate.
func (r *Runtime) ShowBanner() {
	if r.Opts.NoBanner || r.Opts.NonInteractive {
		return
	}
	if !brand.IsTTY(os.Stderr.Fd()) {
		return
	}
	brand.Banner(os.Stderr, r.Opts.Version)
}

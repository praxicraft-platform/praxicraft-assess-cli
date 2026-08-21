package cmdroot

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/brand"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmd/assessments"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmd/cases"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmd/configure"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmd/integrations"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmd/interviews"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmd/invites"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmd/org"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmd/pipelines"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmd/results"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmd/webhooks"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/ui"
	"github.com/spf13/cobra"
)

// Version is injected via ldflags.
var Version = "0.1.3"

// Execute runs the root command and returns a process exit code.
func Execute() int {
	opts := &runtime.Options{Version: Version}
	rt := &runtime.Runtime{}
	*rt = *runtime.New(*opts)

	root := newRoot(opts, rt)
	args := os.Args[1:]

	if len(args) == 0 && brand.IsTTY(os.Stdin.Fd()) && brand.IsTTY(os.Stdout.Fd()) {
		*rt = *runtime.New(*opts)
		rt.Opts.Version = Version
		rt.ShowBanner()
		return runREPL(root, rt)
	}

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return api.ExitCode(err)
	}
	return api.ExitOK
}

func newRoot(opts *runtime.Options, rt *runtime.Runtime) *cobra.Command {
	root := &cobra.Command{
		Use:   "praxicraft-assess",
		Short: "Praxicraft Assess CLI — manage assessments, invites, pipelines, webhooks, and more",
		Long: `Praxicraft Assess CLI — manage assessments, invites, pipelines, webhooks, and more.

Get started:
  praxicraft-assess configure
  praxicraft-assess whoami
  praxicraft-assess assessments list

Docs:      ` + brand.DocsURL + `
CLI guide: ` + brand.CLIDocsURL + `
API keys:  ` + brand.APIKeysURL + `

Run with no arguments in a terminal to open interactive mode (arrow-key resource menu).
Use --help on any command for flags and examples.`,
		Example: `  praxicraft-assess
  praxicraft-assess configure
  praxicraft-assess whoami
  praxicraft-assess assessments list
  praxicraft-assess assessments get
  praxicraft-assess invites create
  praxicraft-assess --query 'results[0].email' invites list
  praxicraft-assess assessments list --filter status=active --output table`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			fresh := runtime.New(*opts)
			fresh.Opts.Version = Version
			*rt = *fresh
		},
	}

	root.PersistentFlags().StringVar(&opts.Profile, "profile", "", "config profile")
	root.PersistentFlags().StringVar(&opts.APIKey, "api-key", "", "API key (or PRAXICRAFT_API_KEY)")
	root.PersistentFlags().StringVar(&opts.BaseURL, "base-url", "", "API base URL")
	root.PersistentFlags().StringVar(&opts.Output, "output", "", "json|table|yaml")
	root.PersistentFlags().StringVar(&opts.Query, "query", "", "JMESPath query on JSON output")
	root.PersistentFlags().BoolVar(&opts.Yes, "yes", false, "skip confirmation prompts")
	root.PersistentFlags().BoolVar(&opts.NonInteractive, "non-interactive", false, "disable prompts")
	root.PersistentFlags().BoolVar(&opts.NoBanner, "no-banner", false, "suppress startup message in interactive mode")

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			brand.VersionLine(os.Stdout, Version)
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "whoami",
		Short: "Show the organisation for the current API key",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := rt.EnsureAPI(); err != nil {
				return err
			}
			out, err := rt.API.OrgGet(rt.Context())
			if err != nil {
				return err
			}
			return rt.Print(out)
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "interactive",
		Short: "Start an interactive Praxicraft Assess shell",
		RunE: func(cmd *cobra.Command, args []string) error {
			*rt = *runtime.New(*opts)
			rt.Opts.Version = Version
			rt.ShowBanner()
			code := runREPL(root, rt)
			if code != 0 {
				os.Exit(code)
			}
			return nil
		},
	})

	configure.Register(root, rt)
	org.Register(root, rt)
	assessments.Register(root, rt)
	invites.Register(root, rt)
	results.Register(root, rt)
	cases.Register(root, rt)
	pipelines.Register(root, rt)
	webhooks.Register(root, rt)
	interviews.Register(root, rt)
	integrations.Register(root, rt)

	return root
}

func runREPL(root *cobra.Command, rt *runtime.Runtime) int {
	fmt.Fprintln(os.Stderr, "Interactive mode — ↑/↓ select a resource command, Enter to run. Type any command, or menu / exit.")
	// Open the picker immediately so the first action is choosing a resource.
	if picked, err := pickREPLCommand(rt); err == nil {
		if code := runREPLLine(root, rt, picked); code >= 0 {
			return code
		}
	}

	for {
		line, err := ui.ReadLine(brand.PromptPrefix())
		if err != nil {
			if ui.IsEOF(err) {
				fmt.Fprintln(os.Stdout)
				return 0
			}
			fmt.Fprintln(os.Stderr, err)
			return api.ExitUsage
		}
		line = strings.TrimSpace(line)
		if line == "" || line == "menu" {
			picked, perr := pickREPLCommand(rt)
			if perr != nil {
				var abort *api.AbortError
				if errors.As(perr, &abort) {
					continue
				}
				fmt.Fprintln(os.Stderr, perr.Error())
				continue
			}
			line = picked
		}
		if code := runREPLLine(root, rt, line); code >= 0 {
			return code
		}
	}
}

// runREPLLine executes one shell line. Returns >=0 to exit the process with that code; -1 to continue.
func runREPLLine(root *cobra.Command, rt *runtime.Runtime, line string) int {
	line = strings.TrimSpace(line)
	if line == "" {
		return -1
	}
	if line == "exit" || line == "quit" {
		return 0
	}
	if line == "help" || line == "--help" || line == "-h" {
		_ = root.Help()
		return -1
	}
	args := splitArgs(line)
	if len(args) > 0 && args[0] == "interactive" {
		fmt.Fprintln(os.Stderr, "already in interactive mode")
		return -1
	}
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	_ = root.Flags().Parse([]string{})
	root.SetArgs(nil)
	return -1
}

func pickREPLCommand(rt *runtime.Runtime) (string, error) {
	choices := []ui.Choice{
		{Label: "whoami — current organisation", Value: "whoami"},
		{Label: "org get — organisation profile", Value: "org get"},
		{Label: "org stats — organisation stats", Value: "org stats"},
		{Label: "org team — team members", Value: "org team"},
		{Label: "assessments list", Value: "assessments list"},
		{Label: "assessments get — pick assessment", Value: "assessments get"},
		{Label: "assessments results — pick assessment", Value: "assessments results"},
		{Label: "assessments cases list — pick assessment", Value: "assessments cases list"},
		{Label: "invites list", Value: "invites list"},
		{Label: "invites create — pick assessment + form", Value: "invites create"},
		{Label: "invites result — pick invite", Value: "invites result"},
		{Label: "invites remind — pick invite", Value: "invites remind"},
		{Label: "invites cancel — pick invite", Value: "invites cancel"},
		{Label: "results list — pick assessment", Value: "results list"},
		{Label: "results get — pick invite", Value: "results get"},
		{Label: "cases list — organisation cases", Value: "cases list"},
		{Label: "cases platform-list — platform catalog", Value: "cases platform-list"},
		{Label: "pipelines list", Value: "pipelines list"},
		{Label: "pipelines get — pick pipeline", Value: "pipelines get"},
		{Label: "pipelines enrollments — pick pipeline", Value: "pipelines enrollments"},
		{Label: "webhooks list", Value: "webhooks list"},
		{Label: "webhooks get — pick webhook", Value: "webhooks get"},
		{Label: "webhooks deliveries — pick webhook", Value: "webhooks deliveries"},
		{Label: "webhooks test — pick webhook", Value: "webhooks test"},
		{Label: "interviews list", Value: "interviews list"},
		{Label: "integrations list", Value: "integrations list"},
		{Label: "configure — API key & profile", Value: "configure"},
		{Label: "help", Value: "help"},
		{Label: "exit", Value: "exit"},
	}
	return ui.Select(rt.UI, "What do you want to do?", choices)
}

func splitArgs(line string) []string {
	var out []string
	var cur strings.Builder
	inQ := false
	for _, r := range line {
		switch {
		case r == '"' && !inQ:
			inQ = true
		case r == '"' && inQ:
			inQ = false
		case r == ' ' && !inQ:
			if cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}

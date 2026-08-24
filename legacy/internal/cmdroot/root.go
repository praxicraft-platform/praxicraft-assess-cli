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
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/config"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/ui"
	"github.com/spf13/cobra"
)

// Version is injected via ldflags.
var Version = "0.1.4"

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
		var abort *api.AbortError
		if errors.As(err, &abort) {
			ui.Aborted()
			return api.ExitOK
		}
		ui.Fail(err)
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
	if picked, err := pickREPLCommand(rt); err == nil {
		if code := runREPLLine(root, rt, picked); code >= 0 {
			return code
		}
	} else {
		var abort *api.AbortError
		if !errors.As(err, &abort) {
			ui.Fail(err)
		}
	}

	for {
		line, err := ui.ReadLine(brand.PromptPrefix())
		if err != nil {
			if ui.IsEOF(err) {
				fmt.Fprintln(os.Stdout)
				return 0
			}
			ui.Fail(err)
			return api.ExitUsage
		}
		line = strings.TrimSpace(line)
		if line == "" || line == "menu" {
			picked, perr := pickREPLCommand(rt)
			if perr != nil {
				var abort *api.AbortError
				if errors.As(perr, &abort) {
					ui.Aborted()
					continue
				}
				ui.Fail(perr)
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
		ui.Done("goodbye")
		return 0
	}
	if line == "help" || line == "--help" || line == "-h" {
		_ = root.Help()
		return -1
	}
	args := splitArgs(line)
	if len(args) > 0 && args[0] == "interactive" {
		ui.Warn("already in interactive mode")
		return -1
	}
	root.SetArgs(args)
	err := root.Execute()
	_ = root.Flags().Parse([]string{})
	root.SetArgs(nil)
	if err != nil {
		var abort *api.AbortError
		if errors.As(err, &abort) {
			ui.Aborted()
			return -1
		}
		ui.Fail(err)
		return -1
	}
	ui.Done("done · Enter opens command center · or type a command")
	return -1
}

func pickREPLCommand(rt *runtime.Runtime) (string, error) {
	act := func(key, label, hint, cmd string) ui.MenuNode {
		return ui.MenuNode{Key: key, Label: label, Hint: hint, Command: cmd}
	}
	menu := []ui.MenuNode{
		{
			Key: "1", Label: "Assessments", Hint: "screens · cases · results",
			Children: []ui.MenuNode{
				act("l", "List", "all assessments", "assessments list"),
				act("g", "Get", "open one (pick)", "assessments get"),
				act("r", "Results", "scores for one assessment", "assessments results"),
				act("c", "Cases", "cases on an assessment", "assessments cases list"),
			},
		},
		{
			Key: "2", Label: "Invites", Hint: "create · remind · cancel",
			Children: []ui.MenuNode{
				act("l", "List", "all invitations", "invites list"),
				act("c", "Create", "invite a candidate", "invites create"),
				act("r", "Result", "outcome for one invite", "invites result"),
				act("m", "Remind", "nudge a candidate", "invites remind"),
				act("x", "Cancel", "revoke an invite", "invites cancel"),
			},
		},
		{
			Key: "3", Label: "Results", Hint: "scores by assessment or token",
			Children: []ui.MenuNode{
				act("l", "List", "by assessment", "results list"),
				act("g", "Get", "by invite token", "results get"),
			},
		},
		{
			Key: "4", Label: "Pipelines", Hint: "stages · enrollments",
			Children: []ui.MenuNode{
				act("l", "List", "all pipelines", "pipelines list"),
				act("g", "Get", "open one (pick)", "pipelines get"),
				act("e", "Enrollments", "candidates in a pipeline", "pipelines enrollments"),
			},
		},
		{
			Key: "5", Label: "Webhooks", Hint: "endpoints · deliveries · test",
			Children: []ui.MenuNode{
				act("l", "List", "all endpoints", "webhooks list"),
				act("g", "Get", "open one (pick)", "webhooks get"),
				act("d", "Deliveries", "delivery log", "webhooks deliveries"),
				act("t", "Test", "send a test event", "webhooks test"),
			},
		},
		{
			Key: "6", Label: "Cases", Hint: "org library · platform catalog",
			Children: []ui.MenuNode{
				act("o", "Organisation", "your case library", "cases list"),
				act("p", "Platform", "shared catalog", "cases platform-list"),
			},
		},
		{
			Key: "7", Label: "Organisation", Hint: "whoami · team · stats",
			Children: []ui.MenuNode{
				act("w", "Whoami", "current organisation", "whoami"),
				act("g", "Profile", "org get", "org get"),
				act("s", "Stats", "invite quota & usage", "org stats"),
				act("t", "Team", "members", "org team"),
			},
		},
		{
			Key: "8", Label: "Interviews", Hint: "live rooms",
			Children: []ui.MenuNode{
				act("l", "List", "interview rooms", "interviews list"),
			},
		},
		{
			Key: "9", Label: "System", Hint: "configure · integrations · help",
			Children: []ui.MenuNode{
				act("c", "Configure", "API key & profile", "configure"),
				act("i", "Integrations", "connected ATS", "integrations list"),
				act("h", "Help", "command help", "help"),
				act("x", "Exit", "leave interactive mode", "exit"),
			},
		},
	}

	session := ui.SessionInfo{
		Version: rt.Opts.Version,
		Profile: rt.Opts.Profile,
		BaseURL: rt.Opts.BaseURL,
	}
	// Best-effort session enrichment (never blocks the menu).
	if res, err := resolveSession(rt); err == nil {
		if session.Profile == "" {
			session.Profile = res.profile
		}
		if session.BaseURL == "" {
			session.BaseURL = res.baseURL
		}
		session.OrgName = res.org
	}

	return ui.Shell(rt.UI, session, menu)
}

type sessionResolved struct {
	profile string
	baseURL string
	org     string
}

func resolveSession(rt *runtime.Runtime) (sessionResolved, error) {
	out := sessionResolved{}
	res, err := config.Resolve(rt.Opts.Profile, rt.Opts.APIKey, rt.Opts.BaseURL)
	if err != nil {
		return out, err
	}
	out.profile = res.Profile
	out.baseURL = res.BaseURL
	if err := rt.EnsureAPI(); err != nil {
		return out, nil
	}
	raw, err := rt.API.OrgGet(rt.Context())
	if err != nil {
		return out, nil
	}
	if m, ok := raw.(map[string]any); ok {
		for _, k := range []string{"name", "organisation_name", "org_name", "slug"} {
			if v, ok := m[k]; ok && v != nil {
				s := strings.TrimSpace(fmt.Sprint(v))
				if s != "" && s != "<nil>" {
					out.org = s
					break
				}
			}
		}
	}
	return out, nil
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

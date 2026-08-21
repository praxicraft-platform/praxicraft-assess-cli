package configure

import (
	"fmt"
	"os"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/brand"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/config"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/ui"
	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	var name string
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Set API key and base URL for a profile",
		Long: `Register this CLI with your Praxicraft Assess organisation.

Create an API key in the dashboard (Assess → Developer → API Keys), then run
configure and paste it when prompted. Keys are shown once (ct_live_… / ct_test_…).

  Product:   ` + brand.ProductURL + `
  Docs:      ` + brand.DocsURL + `
  CLI guide: ` + brand.CLIDocsURL + `
  API keys:  ` + brand.APIKeysURL + `
`,
		RunE: func(c *cobra.Command, args []string) error {
			if name == "" {
				name = "default"
			}
			apiKey := rt.Opts.APIKey
			baseURL := rt.Opts.BaseURL
			var err error

			brand.ConfigureIntro(os.Stdout)

			ui.Section("Authentication")

			if apiKey == "" {
				fmt.Println("Paste a ct_live_… or ct_test_… key from Assess → Developer → API Keys.")
				fmt.Println("Create one here:", brand.APIKeysURL)
				fmt.Println()
				apiKey, err = ui.PromptSecretEnter(rt.UI, "Enter your Praxicraft Assess API key:")
				if err != nil {
					return err
				}
			} else {
				fmt.Println("Using API key from flags or environment.")
			}
			ui.OK("API key accepted")

			ui.Section("Configuration")

			if baseURL == "" {
				baseURL, err = ui.PromptEnter(rt.UI, "Enter the API base URL:", brand.DefaultBaseURL)
				if err != nil {
					return err
				}
			} else {
				fmt.Printf("Using base URL %s\n", baseURL)
			}

			name, err = ui.PromptEnter(rt.UI, "Enter the name of the profile to save:", name)
			if err != nil {
				return err
			}

			f, err := config.Load()
			if err != nil {
				return err
			}
			if f.Profiles == nil {
				f.Profiles = map[string]config.Profile{}
			}
			f.Profiles[name] = config.Profile{APIKey: apiKey, BaseURL: baseURL}
			if f.DefaultProfile == "" {
				f.DefaultProfile = name
			}
			if err := config.Save(f); err != nil {
				return err
			}
			path, _ := config.ConfigPath()
			ui.OK("Settings Saved.")
			fmt.Printf("\nProfile %q written to %s\n", name, path)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "default", "profile name to write")

	list := &cobra.Command{Use: "list", Short: "List configured profiles", RunE: func(c *cobra.Command, args []string) error {
		f, err := config.Load()
		if err != nil {
			return err
		}
		type row struct {
			Name    string `json:"name"`
			BaseURL string `json:"base_url"`
			HasKey  bool   `json:"has_api_key"`
			Default bool   `json:"default"`
		}
		var rows []row
		for n, p := range f.Profiles {
			rows = append(rows, row{Name: n, BaseURL: p.BaseURL, HasKey: p.APIKey != "", Default: n == f.DefaultProfile})
		}
		return rt.Out.Print(rows)
	}}
	cmd.AddCommand(list)
	parent.AddCommand(cmd)
}

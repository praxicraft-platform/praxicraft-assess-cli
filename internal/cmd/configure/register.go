package configure

import (
	"fmt"

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
		RunE: func(c *cobra.Command, args []string) error {
			if name == "" {
				name = "default"
			}
			apiKey := rt.Opts.APIKey
			baseURL := rt.Opts.BaseURL
			var err error
			if apiKey == "" {
				apiKey, err = ui.PromptSecret(rt.UI, "API key (ct_live_… / ct_test_…)")
				if err != nil {
					return err
				}
			}
			if baseURL == "" {
				baseURL, err = ui.PromptString(rt.UI, "Base URL", "https://assess.praxicraft.com")
				if err != nil {
					return err
				}
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
			fmt.Printf("Saved profile %q to %s\n", name, path)
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

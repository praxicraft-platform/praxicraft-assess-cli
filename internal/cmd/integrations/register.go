package integrations

import (
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdutil"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	cmd := &cobra.Command{Use: "integrations", Short: "ATS / provider integrations"}
	cmd.AddCommand(&cobra.Command{Use: "list", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.IntegrationsList(rt.Context()) })
	}})
	cmd.AddCommand(&cobra.Command{Use: "connect-url [provider]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.IntegrationsConnectURL(rt.Context(), args[0]) })
	}})
	cmd.AddCommand(&cobra.Command{Use: "test [provider]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.IntegrationsTest(rt.Context(), args[0]) })
	}})
	parent.AddCommand(cmd)
}

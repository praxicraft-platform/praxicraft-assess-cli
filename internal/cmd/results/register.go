package results

import (
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdutil"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	cmd := &cobra.Command{Use: "results", Short: "Assessment and invite results"}
	var listQ []string
	list := &cobra.Command{Use: "list [assessment-slug]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.ResultsList(rt.Context(), args[0], cmdutil.QueryFromPairs(listQ)) })
	}}
	list.Flags().StringArrayVar(&listQ, "query", nil, "query key=value")
	cmd.AddCommand(list)
	cmd.AddCommand(&cobra.Command{Use: "get [invite-token]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.ResultsGet(rt.Context(), args[0]) })
	}})
	parent.AddCommand(cmd)
}

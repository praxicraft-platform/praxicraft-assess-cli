package results

import (
	"net/url"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdutil"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	cmd := &cobra.Command{Use: "results", Short: "Assessment and invite results"}
	var listQ []string
	var listAll bool
	list := &cobra.Command{
		Use:   "list [assessment-slug]",
		Args:  cobra.MaximumNArgs(1),
		Short: "List results for an assessment (pick interactively if omitted)",
		RunE: func(c *cobra.Command, args []string) error {
			slug := ""
			if len(args) > 0 {
				slug = args[0]
			}
			if slug == "" {
				var err error
				slug, err = cmdutil.PickAssessmentSlug(rt)
				if err != nil {
					return err
				}
			}
			return cmdutil.Run(rt, func() (any, error) {
				return cmdutil.ListOrAll(listAll, listQ, func(q url.Values) (any, error) {
					return rt.API.ResultsList(rt.Context(), slug, q)
				})
			})
		},
	}
	cmdutil.FilterFlag(list, &listQ)
	cmdutil.AllFlag(list, &listAll)
	cmd.AddCommand(list)
	cmd.AddCommand(&cobra.Command{
		Use:   "get [invite-token]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Get result by invite token (pick interactively if omitted)",
		RunE: func(c *cobra.Command, args []string) error {
			token := ""
			if len(args) > 0 {
				token = args[0]
			}
			if token == "" {
				var err error
				token, err = cmdutil.PickInviteToken(rt)
				if err != nil {
					return err
				}
			}
			return cmdutil.Run(rt, func() (any, error) { return rt.API.ResultsGet(rt.Context(), token) })
		},
	})
	parent.AddCommand(cmd)
}

package cases

import (
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdutil"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	cmd := &cobra.Command{Use: "cases", Short: "Platform and organisation cases"}
	var pq, lq []string
	pl := &cobra.Command{Use: "platform-list", Short: "List platform case catalog", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.CasesPlatformList(rt.Context(), cmdutil.QueryFromPairs(pq)) })
	}}
	cmdutil.FilterFlag(pl, &pq)
	cmd.AddCommand(pl)

	list := &cobra.Command{Use: "list", Short: "List organisation cases", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.CasesList(rt.Context(), cmdutil.QueryFromPairs(lq)) })
	}}
	cmdutil.FilterFlag(list, &lq)
	cmd.AddCommand(list)

	var body, file string
	create := &cobra.Command{Use: "create", Short: "Create organisation case", RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(body, file)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.CasesCreate(rt.Context(), b) })
	}}
	cmdutil.BodyFlags(create, &body, &file)
	cmd.AddCommand(create)

	cmd.AddCommand(&cobra.Command{Use: "get [id]", Args: cobra.ExactArgs(1), Short: "Get case by id", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.CasesGet(rt.Context(), args[0]) })
	}})

	var ubody, ufile string
	update := &cobra.Command{Use: "update [id]", Args: cobra.ExactArgs(1), Short: "Update case", RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(ubody, ufile)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.CasesUpdate(rt.Context(), args[0], b) })
	}}
	cmdutil.BodyFlags(update, &ubody, &ufile)
	cmd.AddCommand(update)

	cmd.AddCommand(&cobra.Command{Use: "delete [id]", Args: cobra.ExactArgs(1), Short: "Delete case", RunE: func(c *cobra.Command, args []string) error {
		if err := cmdutil.ConfirmDestructive(rt, "Delete this case?"); err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.CasesDelete(rt.Context(), args[0]) })
	}})
	parent.AddCommand(cmd)
}

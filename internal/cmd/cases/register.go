package cases

import (
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdutil"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	cmd := &cobra.Command{Use: "cases", Short: "Platform and organisation cases"}
	var pq, lq []string
	pl := &cobra.Command{Use: "platform-list", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.CasesPlatformList(rt.Context(), cmdutil.QueryFromPairs(pq)) })
	}}
	pl.Flags().StringArrayVar(&pq, "query", nil, "query key=value")
	cmd.AddCommand(pl)
	list := &cobra.Command{Use: "list", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.CasesList(rt.Context(), cmdutil.QueryFromPairs(lq)) })
	}}
	list.Flags().StringArrayVar(&lq, "query", nil, "query key=value")
	cmd.AddCommand(list)

	var body, file string
	create := &cobra.Command{Use: "create", RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(body, file)
		if err != nil { return err }
		return cmdutil.Run(rt, func() (any, error) { return rt.API.CasesCreate(rt.Context(), b) })
	}}
	cmdutil.BodyFlags(create, &body, &file)
	cmd.AddCommand(create)

	cmd.AddCommand(&cobra.Command{Use: "get [id]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.CasesGet(rt.Context(), args[0]) })
	}})
	var ubody, ufile string
	update := &cobra.Command{Use: "update [id]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(ubody, ufile)
		if err != nil { return err }
		return cmdutil.Run(rt, func() (any, error) { return rt.API.CasesUpdate(rt.Context(), args[0], b) })
	}}
	cmdutil.BodyFlags(update, &ubody, &ufile)
	cmd.AddCommand(update)
	cmd.AddCommand(&cobra.Command{Use: "delete [id]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := cmdutil.ConfirmDestructive(rt, "Delete this case?"); err != nil { return err }
		return cmdutil.Run(rt, func() (any, error) { return rt.API.CasesDelete(rt.Context(), args[0]) })
	}})
	parent.AddCommand(cmd)
}

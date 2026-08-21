package webhooks

import (
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdutil"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	cmd := &cobra.Command{Use: "webhooks", Short: "Webhook endpoints and deliveries"}
	cmd.AddCommand(&cobra.Command{Use: "list", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksList(rt.Context()) })
	}})
	var body, file string
	create := &cobra.Command{Use: "create", RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(body, file)
		if err != nil { return err }
		return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksCreate(rt.Context(), b) })
	}}
	cmdutil.BodyFlags(create, &body, &file)
	cmd.AddCommand(create)
	cmd.AddCommand(&cobra.Command{Use: "get [id]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksGet(rt.Context(), args[0]) })
	}})
	var ubody, ufile string
	update := &cobra.Command{Use: "update [id]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(ubody, ufile)
		if err != nil { return err }
		return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksUpdate(rt.Context(), args[0], b) })
	}}
	cmdutil.BodyFlags(update, &ubody, &ufile)
	cmd.AddCommand(update)
	cmd.AddCommand(&cobra.Command{Use: "delete [id]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := cmdutil.ConfirmDestructive(rt, "Delete this webhook?"); err != nil { return err }
		return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksDelete(rt.Context(), args[0]) })
	}})
	cmd.AddCommand(&cobra.Command{Use: "test [id]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksTest(rt.Context(), args[0]) })
	}})
	var dq []string
	deliv := &cobra.Command{Use: "deliveries [id]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksDeliveries(rt.Context(), args[0], cmdutil.QueryFromPairs(dq)) })
	}}
	deliv.Flags().StringArrayVar(&dq, "query", nil, "query key=value")
	cmd.AddCommand(deliv)
	cmd.AddCommand(&cobra.Command{Use: "retry-delivery [webhook-id] [delivery-id]", Args: cobra.ExactArgs(2), RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksRetryDelivery(rt.Context(), args[0], args[1]) })
	}})
	parent.AddCommand(cmd)
}

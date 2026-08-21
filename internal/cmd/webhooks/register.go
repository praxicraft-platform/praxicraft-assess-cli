package webhooks

import (
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdutil"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/spf13/cobra"
)

func resolveID(rt *runtime.Runtime, args []string) (string, error) {
	if len(args) > 0 && args[0] != "" {
		return args[0], nil
	}
	return cmdutil.PickWebhookID(rt)
}

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	cmd := &cobra.Command{Use: "webhooks", Short: "Webhook endpoints and deliveries"}
	cmd.AddCommand(&cobra.Command{Use: "list", Short: "List webhook endpoints", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksList(rt.Context()) })
	}})

	var body, file string
	create := &cobra.Command{Use: "create", Short: "Create webhook endpoint", RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(body, file)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksCreate(rt.Context(), b) })
	}}
	cmdutil.BodyFlags(create, &body, &file)
	cmd.AddCommand(create)

	cmd.AddCommand(&cobra.Command{
		Use:   "get [id]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Get webhook endpoint (pick interactively if omitted)",
		RunE: func(c *cobra.Command, args []string) error {
			id, err := resolveID(rt, args)
			if err != nil {
				return err
			}
			return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksGet(rt.Context(), id) })
		},
	})

	var ubody, ufile string
	update := &cobra.Command{
		Use:   "update [id]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Update webhook endpoint",
		RunE: func(c *cobra.Command, args []string) error {
			id, err := resolveID(rt, args)
			if err != nil {
				return err
			}
			b, err := cmdutil.ReadBody(ubody, ufile)
			if err != nil {
				return err
			}
			return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksUpdate(rt.Context(), id, b) })
		},
	}
	cmdutil.BodyFlags(update, &ubody, &ufile)
	cmd.AddCommand(update)

	cmd.AddCommand(&cobra.Command{
		Use:   "delete [id]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Delete webhook endpoint",
		RunE: func(c *cobra.Command, args []string) error {
			id, err := resolveID(rt, args)
			if err != nil {
				return err
			}
			if err := cmdutil.ConfirmDestructive(rt, "Delete this webhook?"); err != nil {
				return err
			}
			return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksDelete(rt.Context(), id) })
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "test [id]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Send a test event",
		RunE: func(c *cobra.Command, args []string) error {
			id, err := resolveID(rt, args)
			if err != nil {
				return err
			}
			return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksTest(rt.Context(), id) })
		},
	})

	var dq []string
	deliv := &cobra.Command{
		Use:   "deliveries [id]",
		Args:  cobra.MaximumNArgs(1),
		Short: "List deliveries for a webhook",
		RunE: func(c *cobra.Command, args []string) error {
			id, err := resolveID(rt, args)
			if err != nil {
				return err
			}
			return cmdutil.Run(rt, func() (any, error) {
				return rt.API.WebhooksDeliveries(rt.Context(), id, cmdutil.QueryFromPairs(dq))
			})
		},
	}
	cmdutil.FilterFlag(deliv, &dq)
	cmd.AddCommand(deliv)

	cmd.AddCommand(&cobra.Command{Use: "retry-delivery [webhook-id] [delivery-id]", Args: cobra.ExactArgs(2), Short: "Retry a failed delivery", RunE: func(c *cobra.Command, args []string) error {
		if err := cmdutil.ConfirmDestructive(rt, "Retry this delivery?"); err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.WebhooksRetryDelivery(rt.Context(), args[0], args[1]) })
	}})
	parent.AddCommand(cmd)
}

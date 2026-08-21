package invites

import (
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdutil"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/ui"
	"github.com/spf13/cobra"
)

func resolveAssessment(rt *runtime.Runtime, args []string) (string, error) {
	if len(args) > 0 && args[0] != "" {
		return args[0], nil
	}
	return cmdutil.PickAssessmentSlug(rt)
}

func resolveInvite(rt *runtime.Runtime, args []string) (string, error) {
	if len(args) > 0 && args[0] != "" {
		return args[0], nil
	}
	return cmdutil.PickInviteToken(rt)
}

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	cmd := &cobra.Command{Use: "invites", Short: "Create and manage candidate invitations"}
	var listQ []string
	list := &cobra.Command{Use: "list", Short: "List invitations", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InvitesList(rt.Context(), cmdutil.QueryFromPairs(listQ)) })
	}}
	cmdutil.FilterFlag(list, &listQ)
	cmd.AddCommand(list)

	var email, name string
	var sendEmail bool
	var createBody, createFile string
	create := &cobra.Command{
		Use:   "create [assessment-slug]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Invite a candidate (pick assessment + form if interactive)",
		RunE: func(c *cobra.Command, args []string) error {
			slug, err := resolveAssessment(rt, args)
			if err != nil {
				return err
			}
			body, err := cmdutil.ReadBody(createBody, createFile)
			if err != nil {
				return err
			}
			// Only run the interactive form when flags/body did not already supply email.
			if email == "" && (body == nil || body["email"] == nil || body["email"] == "") {
				email, name, sendEmail, err = ui.FormInvite(rt.UI, email, name, sendEmail)
				if err != nil {
					return err
				}
			} else if email == "" && body != nil {
				if e, ok := body["email"].(string); ok {
					email = e
				}
			}
			if body == nil {
				body = map[string]any{}
			}
			if email != "" {
				body["email"] = email
			}
			if name != "" {
				body["name"] = name
			}
			body["send_email"] = sendEmail
			return cmdutil.Run(rt, func() (any, error) { return rt.API.InvitesCreate(rt.Context(), slug, body) })
		},
	}
	create.Flags().StringVar(&email, "email", "", "candidate email")
	create.Flags().StringVar(&name, "name", "", "candidate name")
	create.Flags().BoolVar(&sendEmail, "send-email", true, "send invite email")
	cmdutil.BodyFlags(create, &createBody, &createFile)
	cmd.AddCommand(create)

	var bulkBody, bulkFile string
	bulk := &cobra.Command{
		Use:   "bulk-create [assessment-slug]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Create many invites from JSON body",
		RunE: func(c *cobra.Command, args []string) error {
			slug, err := resolveAssessment(rt, args)
			if err != nil {
				return err
			}
			body, err := cmdutil.ReadBody(bulkBody, bulkFile)
			if err != nil {
				return err
			}
			return cmdutil.Run(rt, func() (any, error) { return rt.API.InvitesBulkCreate(rt.Context(), slug, body) })
		},
	}
	cmdutil.BodyFlags(bulk, &bulkBody, &bulkFile)
	cmd.AddCommand(bulk)

	cmd.AddCommand(&cobra.Command{
		Use:   "result [invite-token]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Get invite result",
		RunE: func(c *cobra.Command, args []string) error {
			token, err := resolveInvite(rt, args)
			if err != nil {
				return err
			}
			return cmdutil.Run(rt, func() (any, error) { return rt.API.InvitesResult(rt.Context(), token) })
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "remind [invite-token]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Send invite reminder",
		RunE: func(c *cobra.Command, args []string) error {
			token, err := resolveInvite(rt, args)
			if err != nil {
				return err
			}
			if err := cmdutil.ConfirmDestructive(rt, "Send a reminder for this invite?"); err != nil {
				return err
			}
			return cmdutil.Run(rt, func() (any, error) { return rt.API.InvitesRemind(rt.Context(), token) })
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "cancel [invite-token]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Cancel an invite",
		RunE: func(c *cobra.Command, args []string) error {
			token, err := resolveInvite(rt, args)
			if err != nil {
				return err
			}
			if err := cmdutil.ConfirmDestructive(rt, "Cancel this invite?"); err != nil {
				return err
			}
			return cmdutil.Run(rt, func() (any, error) { return rt.API.InvitesCancel(rt.Context(), token) })
		},
	})
	parent.AddCommand(cmd)
}

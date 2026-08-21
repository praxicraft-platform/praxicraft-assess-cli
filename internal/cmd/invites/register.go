package invites

import (
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdutil"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/ui"
	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	cmd := &cobra.Command{Use: "invites", Short: "Create and manage candidate invitations"}
	var listQ []string
	list := &cobra.Command{Use: "list", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InvitesList(rt.Context(), cmdutil.QueryFromPairs(listQ)) })
	}}
	list.Flags().StringArrayVar(&listQ, "query", nil, "query key=value")
	cmd.AddCommand(list)

	var email, name string
	var sendEmail bool
	var createBody, createFile string
	create := &cobra.Command{Use: "create [assessment-slug]", Args: cobra.ExactArgs(1), Short: "Invite a candidate", RunE: func(c *cobra.Command, args []string) error {
		body, err := cmdutil.ReadBody(createBody, createFile)
		if err != nil { return err }
		var e error
		email, e = ui.PromptString(rt.UI, "Candidate email", email)
		if e != nil { return e }
		if body == nil { body = map[string]any{} }
		if email != "" { body["email"] = email }
		if name != "" { body["name"] = name }
		body["send_email"] = sendEmail
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InvitesCreate(rt.Context(), args[0], body) })
	}}
	create.Flags().StringVar(&email, "email", "", "candidate email")
	create.Flags().StringVar(&name, "name", "", "candidate name")
	create.Flags().BoolVar(&sendEmail, "send-email", true, "send invite email")
	cmdutil.BodyFlags(create, &createBody, &createFile)
	cmd.AddCommand(create)

	var bulkBody, bulkFile string
	bulk := &cobra.Command{Use: "bulk-create [assessment-slug]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		body, err := cmdutil.ReadBody(bulkBody, bulkFile)
		if err != nil { return err }
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InvitesBulkCreate(rt.Context(), args[0], body) })
	}}
	cmdutil.BodyFlags(bulk, &bulkBody, &bulkFile)
	cmd.AddCommand(bulk)

	cmd.AddCommand(&cobra.Command{Use: "result [invite-token]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InvitesResult(rt.Context(), args[0]) })
	}})
	cmd.AddCommand(&cobra.Command{Use: "remind [invite-token]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InvitesRemind(rt.Context(), args[0]) })
	}})
	cmd.AddCommand(&cobra.Command{Use: "cancel [invite-token]", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := cmdutil.ConfirmDestructive(rt, "Cancel this invite?"); err != nil { return err }
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InvitesCancel(rt.Context(), args[0]) })
	}})
	parent.AddCommand(cmd)
}

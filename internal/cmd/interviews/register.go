package interviews

import (
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdutil"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	cmd := &cobra.Command{Use: "interviews", Short: "Live interview rooms and templates"}
	var listQ, aq, oq []string
	list := &cobra.Command{Use: "list", Short: "List interview rooms", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewsList(rt.Context(), cmdutil.QueryFromPairs(listQ)) })
	}}
	cmdutil.FilterFlag(list, &listQ)
	cmd.AddCommand(list)

	var cbody, cfile string
	create := &cobra.Command{Use: "create", Short: "Create interview room", RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(cbody, cfile)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewsCreate(rt.Context(), b) })
	}}
	cmdutil.BodyFlags(create, &cbody, &cfile)
	cmd.AddCommand(create)

	var bbody, bfile string
	bulk := &cobra.Command{Use: "bulk-create", Short: "Bulk create interview rooms", RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(bbody, bfile)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewsBulkCreate(rt.Context(), b) })
	}}
	cmdutil.BodyFlags(bulk, &bbody, &bfile)
	cmd.AddCommand(bulk)

	analytics := &cobra.Command{Use: "analytics", Short: "Interview analytics", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewsAnalytics(rt.Context(), cmdutil.QueryFromPairs(aq)) })
	}}
	cmdutil.FilterFlag(analytics, &aq)
	cmd.AddCommand(analytics)

	orgCases := &cobra.Command{Use: "org-cases", Short: "Cases available for interviews", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewsOrgCases(rt.Context(), cmdutil.QueryFromPairs(oq)) })
	}}
	cmdutil.FilterFlag(orgCases, &oq)
	cmd.AddCommand(orgCases)

	cmd.AddCommand(&cobra.Command{Use: "get [room-id]", Args: cobra.ExactArgs(1), Short: "Get interview room", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewsGet(rt.Context(), args[0]) })
	}})

	var cancelBody, cancelFile string
	cancel := &cobra.Command{Use: "cancel [room-id]", Args: cobra.ExactArgs(1), Short: "Cancel interview", RunE: func(c *cobra.Command, args []string) error {
		if err := cmdutil.ConfirmDestructive(rt, "Cancel this interview?"); err != nil {
			return err
		}
		b, err := cmdutil.ReadBody(cancelBody, cancelFile)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewsCancel(rt.Context(), args[0], b) })
	}}
	cmdutil.BodyFlags(cancel, &cancelBody, &cancelFile)
	cmd.AddCommand(cancel)

	var rbody, rfile string
	reschedule := &cobra.Command{Use: "reschedule [room-id]", Args: cobra.ExactArgs(1), Short: "Reschedule interview", RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(rbody, rfile)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewsReschedule(rt.Context(), args[0], b) })
	}}
	cmdutil.BodyFlags(reschedule, &rbody, &rfile)
	cmd.AddCommand(reschedule)

	cmd.AddCommand(&cobra.Command{Use: "analysis [room-id]", Args: cobra.ExactArgs(1), Short: "Get interview analysis", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewsAnalysis(rt.Context(), args[0]) })
	}})
	cmd.AddCommand(&cobra.Command{Use: "replay [room-id]", Args: cobra.ExactArgs(1), Short: "Get interview replay", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewsReplay(rt.Context(), args[0]) })
	}})

	var sbody, sfile string
	share := &cobra.Command{Use: "share [room-id]", Args: cobra.ExactArgs(1), Short: "Share interview", RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(sbody, sfile)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewsShare(rt.Context(), args[0], b) })
	}}
	cmdutil.BodyFlags(share, &sbody, &sfile)
	cmd.AddCommand(share)

	templates := &cobra.Command{Use: "templates", Short: "Interview templates"}
	templates.AddCommand(&cobra.Command{Use: "list", Short: "List templates", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewTemplatesList(rt.Context()) })
	}})
	var tbody, tfile string
	tcreate := &cobra.Command{Use: "create", Short: "Create template", RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(tbody, tfile)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewTemplatesCreate(rt.Context(), b) })
	}}
	cmdutil.BodyFlags(tcreate, &tbody, &tfile)
	templates.AddCommand(tcreate)
	templates.AddCommand(&cobra.Command{Use: "get [id]", Args: cobra.ExactArgs(1), Short: "Get template", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewTemplatesGet(rt.Context(), args[0]) })
	}})
	var tubody, tufile string
	tupdate := &cobra.Command{Use: "update [id]", Args: cobra.ExactArgs(1), Short: "Update template", RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(tubody, tufile)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewTemplatesUpdate(rt.Context(), args[0], b) })
	}}
	cmdutil.BodyFlags(tupdate, &tubody, &tufile)
	templates.AddCommand(tupdate)
	templates.AddCommand(&cobra.Command{Use: "delete [id]", Args: cobra.ExactArgs(1), Short: "Delete template", RunE: func(c *cobra.Command, args []string) error {
		if err := cmdutil.ConfirmDestructive(rt, "Delete this template?"); err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.InterviewTemplatesDelete(rt.Context(), args[0]) })
	}})
	cmd.AddCommand(templates)
	parent.AddCommand(cmd)
}

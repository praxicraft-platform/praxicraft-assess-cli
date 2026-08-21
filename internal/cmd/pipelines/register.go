package pipelines

import (
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdutil"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	cmd := &cobra.Command{Use: "pipelines", Short: "Hiring pipelines and enrollments"}
	var listQ []string
	list := &cobra.Command{Use: "list", Short: "List pipelines", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.PipelinesList(rt.Context(), cmdutil.QueryFromPairs(listQ)) })
	}}
	cmdutil.FilterFlag(list, &listQ)
	cmd.AddCommand(list)

	cmd.AddCommand(&cobra.Command{Use: "get [slug]", Args: cobra.ExactArgs(1), Short: "Get pipeline by slug", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.PipelinesGet(rt.Context(), args[0]) })
	}})

	var ebody, efile string
	enroll := &cobra.Command{Use: "enroll [slug]", Args: cobra.ExactArgs(1), Short: "Enroll a candidate", RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(ebody, efile)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.PipelinesEnroll(rt.Context(), args[0], b) })
	}}
	cmdutil.BodyFlags(enroll, &ebody, &efile)
	cmd.AddCommand(enroll)

	var bbody, bfile string
	bulk := &cobra.Command{Use: "bulk-enroll [slug]", Args: cobra.ExactArgs(1), Short: "Bulk enroll candidates", RunE: func(c *cobra.Command, args []string) error {
		b, err := cmdutil.ReadBody(bbody, bfile)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.PipelinesBulkEnroll(rt.Context(), args[0], b) })
	}}
	cmdutil.BodyFlags(bulk, &bbody, &bfile)
	cmd.AddCommand(bulk)

	var eq []string
	enrollments := &cobra.Command{Use: "enrollments [slug]", Args: cobra.ExactArgs(1), Short: "List enrollments", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.PipelinesEnrollments(rt.Context(), args[0], cmdutil.QueryFromPairs(eq)) })
	}}
	cmdutil.FilterFlag(enrollments, &eq)
	cmd.AddCommand(enrollments)

	cmd.AddCommand(&cobra.Command{Use: "get-enrollment [id]", Args: cobra.ExactArgs(1), Short: "Get enrollment by id", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.PipelinesGetEnrollment(rt.Context(), args[0]) })
	}})

	for _, op := range []struct{ name, short string }{
		{"reject", "Reject an enrollment"},
		{"hold", "Hold an enrollment"},
		{"unhold", "Release an enrollment hold"},
	} {
		op := op
		var body, file string
		cc := &cobra.Command{Use: op.name + " [enrollment-id]", Args: cobra.ExactArgs(1), Short: op.short, RunE: func(c *cobra.Command, args []string) error {
			if err := cmdutil.ConfirmDestructive(rt, op.name+" this enrollment?"); err != nil {
				return err
			}
			b, err := cmdutil.ReadBody(body, file)
			if err != nil {
				return err
			}
			return cmdutil.Run(rt, func() (any, error) {
				switch op.name {
				case "reject":
					return rt.API.PipelinesReject(rt.Context(), args[0], b)
				case "hold":
					return rt.API.PipelinesHold(rt.Context(), args[0], b)
				default:
					return rt.API.PipelinesUnhold(rt.Context(), args[0], b)
				}
			})
		}}
		cmdutil.BodyFlags(cc, &body, &file)
		cmd.AddCommand(cc)
	}
	parent.AddCommand(cmd)
}

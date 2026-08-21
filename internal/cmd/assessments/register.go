package assessments

import (
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdutil"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	cmd := &cobra.Command{Use: "assessments", Short: "Manage assessments and attached cases"}
	var listQ []string
	list := &cobra.Command{Use: "list", Short: "List assessments", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.AssessmentsList(rt.Context(), cmdutil.QueryFromPairs(listQ)) })
	}}
	cmdutil.FilterFlag(list, &listQ)
	cmd.AddCommand(list)

	cmd.AddCommand(&cobra.Command{Use: "get [slug]", Args: cobra.ExactArgs(1), Short: "Get assessment by slug", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.AssessmentsGet(rt.Context(), args[0]) })
	}})

	var createBody, createFile string
	create := &cobra.Command{Use: "create", Short: "Create assessment", RunE: func(c *cobra.Command, args []string) error {
		body, err := cmdutil.ReadBody(createBody, createFile)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.AssessmentsCreate(rt.Context(), body) })
	}}
	cmdutil.BodyFlags(create, &createBody, &createFile)
	cmd.AddCommand(create)

	var updateBody, updateFile string
	update := &cobra.Command{Use: "update [slug]", Args: cobra.ExactArgs(1), Short: "Update assessment", RunE: func(c *cobra.Command, args []string) error {
		body, err := cmdutil.ReadBody(updateBody, updateFile)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.AssessmentsUpdate(rt.Context(), args[0], body) })
	}}
	cmdutil.BodyFlags(update, &updateBody, &updateFile)
	cmd.AddCommand(update)

	var dupBody, dupFile string
	dup := &cobra.Command{Use: "duplicate [slug]", Args: cobra.ExactArgs(1), Short: "Duplicate assessment", RunE: func(c *cobra.Command, args []string) error {
		body, err := cmdutil.ReadBody(dupBody, dupFile)
		if err != nil {
			return err
		}
		return cmdutil.Run(rt, func() (any, error) { return rt.API.AssessmentsDuplicate(rt.Context(), args[0], body) })
	}}
	cmdutil.BodyFlags(dup, &dupBody, &dupFile)
	cmd.AddCommand(dup)

	cases := &cobra.Command{Use: "cases", Short: "Assessment cases"}
	cases.AddCommand(&cobra.Command{Use: "list [slug]", Args: cobra.ExactArgs(1), Short: "List cases on an assessment", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.AssessmentsCasesList(rt.Context(), args[0]) })
	}})
	for _, op := range []struct{ use, method, short string }{
		{"attach [slug]", "attach", "Attach cases to an assessment"},
		{"replace [slug]", "replace", "Replace all cases on an assessment"},
		{"remove [slug]", "remove", "Remove cases from an assessment"},
	} {
		op := op
		var body, file string
		cc := &cobra.Command{Use: op.use, Args: cobra.ExactArgs(1), Short: op.short, RunE: func(c *cobra.Command, args []string) error {
			if op.method == "remove" || op.method == "replace" {
				if err := cmdutil.ConfirmDestructive(rt, op.method+" cases on this assessment?"); err != nil {
					return err
				}
			}
			b, err := cmdutil.ReadBody(body, file)
			if err != nil {
				return err
			}
			return cmdutil.Run(rt, func() (any, error) {
				switch op.method {
				case "attach":
					return rt.API.AssessmentsCasesAttach(rt.Context(), args[0], b)
				case "replace":
					return rt.API.AssessmentsCasesReplace(rt.Context(), args[0], b)
				default:
					return rt.API.AssessmentsCasesRemove(rt.Context(), args[0], b)
				}
			})
		}}
		cmdutil.BodyFlags(cc, &body, &file)
		cases.AddCommand(cc)
	}
	cmd.AddCommand(cases)

	var resQ []string
	res := &cobra.Command{Use: "results [slug]", Args: cobra.ExactArgs(1), Short: "List assessment results", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.AssessmentsResults(rt.Context(), args[0], cmdutil.QueryFromPairs(resQ)) })
	}}
	cmdutil.FilterFlag(res, &resQ)
	cmd.AddCommand(res)
	parent.AddCommand(cmd)
}

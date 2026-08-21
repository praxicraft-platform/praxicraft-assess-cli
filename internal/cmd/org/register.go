package org

import (
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdutil"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command, rt *runtime.Runtime) {
	cmd := &cobra.Command{Use: "org", Short: "Organisation profile, team, stats, and squads"}
	cmd.AddCommand(
		&cobra.Command{Use: "get", Short: "Get organisation profile", RunE: func(c *cobra.Command, args []string) error {
			return cmdutil.Run(rt, func() (any, error) { return rt.API.OrgGet(rt.Context()) })
		}},
		&cobra.Command{Use: "stats", Short: "Get organisation stats / invite quota", RunE: func(c *cobra.Command, args []string) error {
			return cmdutil.Run(rt, func() (any, error) { return rt.API.OrgStats(rt.Context()) })
		}},
		&cobra.Command{Use: "team", Short: "List organisation team members", RunE: func(c *cobra.Command, args []string) error {
			return cmdutil.Run(rt, func() (any, error) { return rt.API.OrgTeam(rt.Context()) })
		}},
	)
	var auditQ []string
	audit := &cobra.Command{Use: "audit-log", Short: "Organisation audit log", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.OrgAuditLog(rt.Context(), cmdutil.QueryFromPairs(auditQ)) })
	}}
	cmdutil.FilterFlag(audit, &auditQ)
	cmd.AddCommand(audit)

	squads := &cobra.Command{Use: "squads", Short: "Squads"}
	squads.AddCommand(&cobra.Command{Use: "list", Short: "List squads", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.OrgSquadsList(rt.Context()) })
	}})
	squads.AddCommand(&cobra.Command{Use: "get [team-id]", Args: cobra.ExactArgs(1), Short: "Get squad", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.OrgSquadGet(rt.Context(), args[0]) })
	}})
	squads.AddCommand(&cobra.Command{Use: "members [team-id]", Args: cobra.ExactArgs(1), Short: "List squad members", RunE: func(c *cobra.Command, args []string) error {
		return cmdutil.Run(rt, func() (any, error) { return rt.API.OrgSquadMembers(rt.Context(), args[0]) })
	}})
	cmd.AddCommand(squads)
	parent.AddCommand(cmd)
}

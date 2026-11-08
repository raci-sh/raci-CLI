package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var realm, dir string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Start a realm's RACI log in this repository",
		Long: "Initialise a realm: append a pin event for each dependency realm (e.g. sh.raci) and the\n" +
			"  rule events of the rules in force, then write an empty .raci view (SPEC.md section 6).",
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if realm == "" {
				return withExitCode(exitUsage, errors.New("--realm is required"))
			}
			return withExitCode(exitNotImplemented, fmt.Errorf("raci init --realm %s: %w", realm, errNotImplemented))
		},
	}
	cmd.Flags().StringVar(&realm, "realm", "", "realm id, matching a SteerSpec realm (e.g. com.example)")
	cmd.Flags().StringVar(&dir, "dir", ".", "directory holding the .raci file and its log")
	return cmd
}

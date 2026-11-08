package cmd

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/raci-sh/raci-manager/src/eventlog"
	"github.com/raci-sh/raci-manager/src/model"
	"github.com/raci-sh/raci-manager/src/verify"
	"github.com/spf13/cobra"
)

func newVerifyCmd() *cobra.Command {
	var realm, dir string
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify a realm's log and .raci view",
		Long: "Recompute the hash chain of the realm's log and print a JSON verification report\n" +
			"  (SPEC.md section 7.1). Exits 1 when the report has errors. Replay, rule and change checks\n" +
			"  are not implemented yet.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			log := eventlog.Open(filepath.Join(dir, eventlog.FileName))
			report, err := verify.Log(model.RealmID(realm), log)
			if err != nil {
				return withExitCode(exitUsage, err)
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			if err := enc.Encode(report); err != nil {
				return err
			}
			if report.HasErrors() {
				return withExitCode(exitFindings, errors.New("verification found errors"))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&realm, "realm", "", "realm id to report on")
	cmd.Flags().StringVar(&dir, "dir", ".", "directory holding the .raci file and its log")
	return cmd
}

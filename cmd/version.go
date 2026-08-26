package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/version"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the agentd build version",
	Long: `Print the build version of this agentd binary.

goreleaser injects a semver via ldflags. go install @tag/@latest uses the
module version from BuildInfo. Local devel builds may show "dev" or
"dev+<shortrev>". This does not contact the daemon. The running service
version is on "agentd daemon status".`,
	Example: `  agentd version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		_, err := fmt.Fprintln(cmd.OutOrStdout(), version.String())
		return err
	},
}

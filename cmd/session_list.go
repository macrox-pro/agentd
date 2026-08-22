package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/provider"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

var (
	sessionListProvider string
	sessionListJSON     bool
)

func init() {
	sessionListCmd.Flags().StringVar(&sessionListProvider, "provider", "", "filter by provider id (claude-code, cursor, …)")
	sessionListCmd.Flags().BoolVar(&sessionListJSON, "json", false, "print JSON")
}

var sessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recorded sessions",
	Long:  `List session JSONL files under the agentd sessions directory.`,
	Example: `  agentd session list
  agentd session list --provider=cursor --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		prov, err := provider.ParseFilter(sessionListProvider, cmd.Flags().Changed("provider"))
		if err != nil {
			return err
		}
		summaries, err := trajectory.ListSessions(trajectory.DefaultSessionsDir(), string(prov))
		if err != nil {
			return err
		}
		if sessionListJSON {
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(summaries)
		}
		if len(summaries) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "no sessions")
			return nil
		}
		for _, s := range summaries {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", s.Provider, s.SessionID)
		}
		return nil
	},
}

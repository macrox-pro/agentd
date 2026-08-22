package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/provider"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

var (
	sessionShowProvider string
	sessionShowJSON     bool
)

func init() {
	sessionShowCmd.Flags().StringVar(&sessionShowProvider, "provider", "", "provider id (required)")
	sessionShowCmd.Flags().BoolVar(&sessionShowJSON, "json", false, "print JSON array of events")
	_ = sessionShowCmd.MarkFlagRequired("provider")
}

var sessionShowCmd = &cobra.Command{
	Use:   "show SESSION_ID",
	Short: "Show events for one session",
	Long:  `Print all events from a session JSONL ledger.`,
	Example: `  agentd session show s1 --provider=claude-code
  agentd session show s1 --provider=claude-code --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prov, err := provider.Parse(sessionShowProvider)
		if err != nil {
			return err
		}
		path, err := trajectory.FindSessionPath(trajectory.DefaultSessionsDir(), string(prov), args[0])
		if err != nil {
			return err
		}
		events, err := trajectory.ReadEvents(path)
		if err != nil {
			return err
		}
		if sessionShowJSON {
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(events)
		}
		for _, e := range events {
			fmt.Fprintf(cmd.OutOrStdout(), "%d\t%s\t%s\n", e.Seq, e.Type, e.Source)
		}
		return nil
	},
}

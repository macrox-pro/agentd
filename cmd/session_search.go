package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

var (
	sessionSearchProvider string
	sessionSearchSession  string
	sessionSearchKinds    []string
	sessionSearchSource   string
	sessionSearchQuery    string
	sessionSearchJSON     bool
	sessionSearchLimit    int
)

func init() {
	sessionSearchCmd.Flags().StringVar(&sessionSearchProvider, "provider", "", "filter by provider id")
	sessionSearchCmd.Flags().StringVar(&sessionSearchSession, "session", "", "filter by session id")
	sessionSearchCmd.Flags().StringSliceVar(&sessionSearchKinds, "kind", nil, "filter by event type (repeatable)")
	sessionSearchCmd.Flags().StringVar(&sessionSearchSource, "source", "", "filter by source (hook, transcript, decision, system)")
	sessionSearchCmd.Flags().StringVar(&sessionSearchQuery, "query", "", "case-insensitive substring match")
	sessionSearchCmd.Flags().BoolVar(&sessionSearchJSON, "json", false, "print JSON array of hits")
	sessionSearchCmd.Flags().IntVar(&sessionSearchLimit, "limit", 100, "maximum hits to return")
}

var sessionSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search session ledger events",
	Long: `Search local session JSONL ledgers by provider, session, kind, source, or text.

Scans every matching file line-by-line (O(total bytes); no search index).`,
	Example: `  agentd session search --query "PreToolUse"
  agentd session search --provider claude-code --source transcript --json
  agentd session search --kind transcript/thinking --limit 20`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		hits, err := trajectory.Search(trajectory.SearchOptions{
			Root:      trajectory.DefaultSessionsDir(),
			Provider:  sessionSearchProvider,
			SessionID: sessionSearchSession,
			Types:     sessionSearchKinds,
			Source:    sessionSearchSource,
			Query:     sessionSearchQuery,
			Limit:     sessionSearchLimit,
		})
		if err != nil {
			return err
		}
		if sessionSearchJSON {
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(hits)
		}
		if len(hits) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "no matches")
			return nil
		}
		for _, h := range hits {
			fmt.Fprintf(cmd.OutOrStdout(), "%s/%s\t%d\t%s\t%s\n", h.Provider, h.SessionID, h.Seq, h.Type, h.Snippet)
		}
		return nil
	},
}

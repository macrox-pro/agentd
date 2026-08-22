package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/trajectory"
)

var (
	sessionForkProvider   string
	sessionForkSession    string
	sessionForkNewSession string
	sessionForkAtSeq      uint64
	sessionForkJSON       bool
)

func init() {
	sessionForkCmd.Flags().StringVar(&sessionForkProvider, "provider", "", "provider id (required)")
	sessionForkCmd.Flags().StringVar(&sessionForkSession, "session", "", "source session id (required)")
	sessionForkCmd.Flags().StringVar(&sessionForkNewSession, "new-session", "", "destination session id (required)")
	sessionForkCmd.Flags().Uint64Var(&sessionForkAtSeq, "at-seq", 0, "copy events with seq <= N (0 = all)")
	sessionForkCmd.Flags().BoolVar(&sessionForkJSON, "json", false, "print result as JSON")
	_ = sessionForkCmd.MarkFlagRequired("provider")
	_ = sessionForkCmd.MarkFlagRequired("session")
	_ = sessionForkCmd.MarkFlagRequired("new-session")
}

var sessionForkCmd = &cobra.Command{
	Use:   "fork",
	Short: "Copy a session ledger prefix into a new session id",
	Long: `Create a new append-only ledger from a prefix of an existing session (audit lineage).

Does not spawn or resume a live agent. The source session JSONL is left immutable.`,
	Example: `  agentd session fork --provider claude-code --session s1 --new-session s1-fork
  agentd session fork --provider cursor --session s1 --new-session s2 --at-seq 4 --json`,
	RunE: runSessionFork,
}

func runSessionFork(cmd *cobra.Command, _ []string) error {
	key := trajectory.ResolveSessionKey(sessionForkProvider, sessionForkSession, "", "")
	result, err := trajectory.ForkSession(trajectory.DefaultSessionsDir(), key, sessionForkNewSession, sessionForkAtSeq)
	if err != nil {
		return err
	}
	if sessionForkJSON {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "provider=%s parent=%s new=%s boundary_seq=%d copied=%d\n",
		result.Provider, result.ParentSession, result.NewSessionID, result.BoundarySeq, result.Copied)
	return nil
}

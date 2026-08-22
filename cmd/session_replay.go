package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/provider"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

var (
	sessionReplayProvider string
	sessionReplaySession  string
	sessionReplaySeq      uint64
	sessionReplayJSON     bool
	sessionReplayPolicy   bool
)

func init() {
	sessionReplayCmd.Flags().BoolVar(&sessionReplayPolicy, "policy", false, "dry-run stored raw through the dispatch engine (required)")
	sessionReplayCmd.Flags().StringVar(&sessionReplayProvider, "provider", "", "provider id (required)")
	sessionReplayCmd.Flags().StringVar(&sessionReplaySession, "session", "", "session id (required)")
	sessionReplayCmd.Flags().Uint64Var(&sessionReplaySeq, "seq", 0, "replay a single hook/invoked seq (0 = all with raw)")
	sessionReplayCmd.Flags().BoolVar(&sessionReplayJSON, "json", false, "print result as JSON")
	_ = sessionReplayCmd.MarkFlagRequired("provider")
	_ = sessionReplayCmd.MarkFlagRequired("session")
}

var sessionReplayCmd = &cobra.Command{
	Use:   "replay",
	Short: "Dry-run policy against stored session payloads",
	Long: `Re-invoke stored hook Raw payloads through the Dispatch Engine (offline).

Requires trajectory.include_raw=true at record time. Does not talk to a live agent
or resume an agent loop — policy dry-run only.`,
	Example: `  agentd session replay --policy --provider claude-code --session s1
  agentd session replay --policy --provider cursor --session s1 --seq 2 --json`,
	RunE: runSessionReplay,
}

func runSessionReplay(cmd *cobra.Command, _ []string) error {
	if !sessionReplayPolicy {
		return fmt.Errorf("--policy is required (agent-loop resume is out of scope)")
	}
	prov, err := provider.Parse(sessionReplayProvider)
	if err != nil {
		return err
	}
	result, err := trajectory.ReplayPolicyFromConfig(cmd.Context(), trajectory.ReplayPolicyConfigOptions{
		ConfigPath:   resolveConfigPath(),
		SessionsRoot: trajectory.DefaultSessionsDir(),
		Provider:     string(prov),
		SessionID:    sessionReplaySession,
		Seq:          sessionReplaySeq,
	})
	if err != nil {
		if errors.Is(err, trajectory.ErrReplayNoRaw) {
			return fmt.Errorf("policy replay requires trajectory.include_raw=true at record time")
		}
		return err
	}

	if sessionReplayJSON {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "provider=%s session=%s hits=%d\n", result.Provider, result.SessionID, len(result.Hits))
	for _, h := range result.Hits {
		if h.Error != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "  seq=%d kind=%s error=%s\n", h.Seq, h.Kind, h.Error)
			continue
		}
		fmt.Fprintf(cmd.OutOrStdout(), "  seq=%d kind=%s stored=%s replay=%s match=%v\n",
			h.Seq, h.Kind, h.StoredDecision, h.ReplayDecision, h.Match)
	}
	return nil
}

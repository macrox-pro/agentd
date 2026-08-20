package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/hookclient"
)

var (
	recordFingerprint string
	recordScope       string
	recordProjectRoot string
	recordSessionID   string
	recordExpiresAt   string
)

func init() {
	configRecordDecisionCmd.Flags().StringVar(&recordFingerprint, "fingerprint", "", "approval fingerprint from Ask system_message (required)")
	configRecordDecisionCmd.Flags().StringVar(&recordScope, "scope", "project", "approval scope: project or session")
	configRecordDecisionCmd.Flags().StringVar(&recordProjectRoot, "project-root", "", "project root for project-scoped approvals")
	configRecordDecisionCmd.Flags().StringVar(&recordSessionID, "session-id", "", "session id for session-scoped approvals")
	configRecordDecisionCmd.Flags().StringVar(&recordExpiresAt, "expires-at", "", "RFC3339 expiry (default: project 24h; session none)")
	_ = configRecordDecisionCmd.MarkFlagRequired("fingerprint")
}

var configRecordDecisionCmd = &cobra.Command{
	Use:   "record-decision",
	Short: "Record an approval after an Ask decision",
	Long: `Record a runtime approval so matching tool.pre calls skip re-Ask within TTL.

Use the approval_fingerprint value from the Ask system_message. Project scope
defaults to a 24h expiry; session scope matches the given session id until cleared.

Requires a running agentd service. Approvals persist to runtime.yaml.`,
	Example: `  agentd config record-decision --fingerprint sha256:secrets/abc... --scope project --project-root /path/to/repo
  agentd config record-decision --fingerprint sha256:shell/def... --scope session --session-id s1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		scope, err := parseRecordScope(recordScope)
		if err != nil {
			return err
		}
		req := &agentdv1.RecordDecisionRequest{
			ApprovalFingerprint: recordFingerprint,
			Scope:               scope,
			ProjectRoot:         recordProjectRoot,
			SessionId:           recordSessionID,
		}
		if recordExpiresAt != "" {
			t, err := time.Parse(time.RFC3339, recordExpiresAt)
			if err != nil {
				return fmt.Errorf("expires-at: %w", err)
			}
			req.ExpiresAt = timestamppb.New(t.UTC())
		}

		cli, err := hookclient.Dial(cmd.Context(), resolveSocket())
		if err != nil {
			return fmt.Errorf("daemon not running: %w", err)
		}
		defer cli.Close()

		resp, err := cli.RecordDecision(cmd.Context(), req)
		if err != nil {
			return err
		}
		cfg := resp.GetConfig()
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "generation=%d fingerprint=%s\n",
			cfg.GetGeneration(), cfg.GetFingerprint())
		return err
	},
}

func parseRecordScope(s string) (agentdv1.ConfigLayer, error) {
	switch s {
	case "project":
		return agentdv1.ConfigLayer_CONFIG_LAYER_PROJECT, nil
	case "session":
		return agentdv1.ConfigLayer_CONFIG_LAYER_RUNTIME, nil
	default:
		return agentdv1.ConfigLayer_CONFIG_LAYER_UNSPECIFIED, fmt.Errorf("scope must be project or session, got %q", s)
	}
}

package guard

import (
	"time"

	"github.com/macrox-pro/agentd/internal/config"
)

// DecisionContext carries compiled approvals for Ask short-circuit.
type DecisionContext struct {
	Approvals       config.Approvals
	TemporaryBlocks []config.TemporaryBlock
	ProjectRoot     string
	AskFallback     config.AskFallback
}

func approved(ctx DecisionContext, kind config.ApprovalKind, fingerprint, sessionID string) bool {
	return ctx.Approvals.HasApproval(kind, fingerprint, ctx.ProjectRoot, sessionID, time.Now().UTC())
}

func findingRuleIDs(findings []Finding) []string {
	ids := make([]string, 0, len(findings))
	for _, f := range findings {
		if f.ID != "" {
			ids = append(ids, f.ID)
		}
	}
	return ids
}

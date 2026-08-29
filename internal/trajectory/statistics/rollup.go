package statistics

import (
	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
)

// StatisticsRollup is daemon-lifetime counter totals.
type StatisticsRollup struct {
	HooksByKind       map[agentdv1.EventKind]uint64
	DecisionsByKind   map[agentdv1.DecisionKind]uint64
	AsyncDispatched   uint64
	InputTokensTotal  uint64
	OutputTokensTotal uint64
	CacheReadTokens   uint64
	CacheWriteTokens  uint64
	ContextTokensLast uint64
}

func newStatisticsRollup() StatisticsRollup {
	return StatisticsRollup{
		HooksByKind:     map[agentdv1.EventKind]uint64{},
		DecisionsByKind: map[agentdv1.DecisionKind]uint64{},
	}
}

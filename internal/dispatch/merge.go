package dispatch

import "github.com/speakeasy-api/agenthooks"

// FirstConclusive reports whether d ends the sync list under first_conclusive
// (non-nil and Kind != DecisionNoDecision). List fold + short-circuit stay in runSync.
func FirstConclusive(d agenthooks.Decision) bool {
	return d != nil && d.Kind() != agenthooks.DecisionNoDecision
}

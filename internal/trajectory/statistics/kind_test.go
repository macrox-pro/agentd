package statistics_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics"
)

func TestHookKind(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		want agentdv1.EventKind
	}{
		{name: "tool_pre", in: "tool.pre", want: agentdv1.EventKind_EVENT_KIND_TOOL_PRE},
		{name: "permission_request", in: "permission.request", want: agentdv1.EventKind_EVENT_KIND_PERMISSION},
		{name: "compact_pre_to_other", in: "compact.pre", want: agentdv1.EventKind_EVENT_KIND_OTHER},
		{name: "subagent_stop_to_other", in: "subagent.stop", want: agentdv1.EventKind_EVENT_KIND_OTHER},
		{name: "subagent_start_to_other", in: "subagent.start", want: agentdv1.EventKind_EVENT_KIND_OTHER},
		{name: "model_response_to_other", in: "model.response", want: agentdv1.EventKind_EVENT_KIND_OTHER},
		{name: "file_edited_to_other", in: "file.edited", want: agentdv1.EventKind_EVENT_KIND_OTHER},
		{name: "agent_stop", in: "agent.stop", want: agentdv1.EventKind_EVENT_KIND_AGENT_STOP},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, statistics.HookKind(tt.in))
		})
	}
}

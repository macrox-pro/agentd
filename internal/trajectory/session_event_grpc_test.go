package trajectory_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestEventSessionEventRoundTrip(t *testing.T) {
	t.Parallel()
	ts := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)
	largeRaw := json.RawMessage(`"` + strings.Repeat("x", 64*1024) + `"`)

	t.Run("full_roundtrip", func(t *testing.T) {
		t.Parallel()
		in := trajectory.Event{
			SchemaVersion:  trajectory.SchemaVersion,
			Seq:            3,
			Type:           trajectory.TypeHookInvoked,
			Source:         trajectory.SourceHook,
			TS:             ts,
			Provider:       "claude-code",
			InvocationMode: "stdin",
			SessionID:      "s1",
			ProjectRoot:    "/proj",
			CWD:            "/proj",
			Data:           json.RawMessage(`{"kind":"tool.pre"}`),
			Raw:            json.RawMessage(`{"hook":"x"}`),
			Ignorable:      true,
		}
		out := trajectory.EventFromSessionEvent(trajectory.EventToSessionEvent(in))
		assert.Equal(t, in.SchemaVersion, out.SchemaVersion)
		assert.Equal(t, in.Seq, out.Seq)
		assert.Equal(t, in.Type, out.Type)
		assert.Equal(t, in.Source, out.Source)
		assert.True(t, in.TS.Equal(out.TS))
		assert.Equal(t, in.Provider, out.Provider)
		assert.Equal(t, in.InvocationMode, out.InvocationMode)
		assert.Equal(t, in.SessionID, out.SessionID)
		assert.Equal(t, in.ProjectRoot, out.ProjectRoot)
		assert.Equal(t, in.CWD, out.CWD)
		assert.JSONEq(t, string(in.Data), string(out.Data))
		assert.JSONEq(t, string(in.Raw), string(out.Raw))
		assert.Equal(t, in.Ignorable, out.Ignorable)
	})

	tests := []struct {
		name  string
		run   func(t *testing.T)
	}{
		{
			name: "nil_proto",
			run: func(t *testing.T) {
				assert.Equal(t, trajectory.Event{}, trajectory.EventFromSessionEvent(nil))
			},
		},
		{
			name: "zero_timestamp",
			run: func(t *testing.T) {
				proto := trajectory.EventToSessionEvent(trajectory.Event{Seq: 1})
				assert.Nil(t, proto.GetTs())
				out := trajectory.EventFromSessionEvent(proto)
				assert.Equal(t, uint64(1), out.Seq)
				assert.True(t, out.TS.IsZero())
			},
		},
		{
			name: "timestamp_encode",
			run: func(t *testing.T) {
				ts := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
				proto := trajectory.EventToSessionEvent(trajectory.Event{TS: ts})
				require.NotNil(t, proto.GetTs())
				assert.True(t, ts.Equal(proto.GetTs().AsTime().UTC()))
			},
		},
		{
			name: "empty_data_raw",
			run: func(t *testing.T) {
				in := trajectory.Event{Seq: 2, Type: trajectory.TypeHookInvoked, Source: trajectory.SourceHook, Provider: "cursor"}
				out := trajectory.EventFromSessionEvent(trajectory.EventToSessionEvent(in))
				assert.Equal(t, in.Seq, out.Seq)
				assert.Empty(t, out.Data)
				assert.Empty(t, out.Raw)
			},
		},
		{
			name: "large_raw",
			run: func(t *testing.T) {
				in := trajectory.Event{
					Seq:      4,
					Type:     trajectory.TypeHookInvoked,
					Source:   trajectory.SourceHook,
					Provider: "codex",
					Raw:      largeRaw,
				}
				out := trajectory.EventFromSessionEvent(trajectory.EventToSessionEvent(in))
				assert.Equal(t, string(in.Raw), string(out.Raw))
			},
		},
		{
			name: "proto_to_event",
			run: func(t *testing.T) {
				proto := &agentdv1.SessionEvent{
					SchemaVersion:  1,
					Seq:            7,
					Type:           trajectory.TypeTranscriptMessage,
					Source:         trajectory.SourceTranscript,
					Ts:             timestamppb.New(time.Unix(100, 0).UTC()),
					Provider:       "cursor",
					InvocationMode: "argv",
					SessionId:      "abc",
					ProjectRoot:    "/p",
					Cwd:            "/p",
					Data:           []byte(`{"text":"hi"}`),
					Raw:            []byte(`{}`),
				}
				ev := trajectory.EventFromSessionEvent(proto)
				assert.Equal(t, uint64(7), ev.Seq)
				assert.Equal(t, "cursor", ev.Provider)
				assert.JSONEq(t, `{"text":"hi"}`, string(ev.Data))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.run(t)
		})
	}
}

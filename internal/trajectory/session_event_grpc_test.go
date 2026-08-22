package trajectory_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestSessionEventGRPCRoundTrip(t *testing.T) {
	t.Parallel()
	ts := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)
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

	proto := trajectory.EventToSessionEvent(in)
	out := trajectory.EventFromSessionEvent(proto)

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
}

func TestEventFromSessionEventNil(t *testing.T) {
	t.Parallel()
	assert.Equal(t, trajectory.Event{}, trajectory.EventFromSessionEvent(nil))
}

func TestEventToSessionEventTimestamp(t *testing.T) {
	t.Parallel()
	ts := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	ev := trajectory.Event{TS: ts}
	proto := trajectory.EventToSessionEvent(ev)
	require.NotNil(t, proto.GetTs())
	assert.True(t, ts.Equal(proto.GetTs().AsTime().UTC()))
}

func TestEventToSessionEventOmitsZeroTimestamp(t *testing.T) {
	t.Parallel()
	proto := trajectory.EventToSessionEvent(trajectory.Event{Seq: 1})
	assert.Nil(t, proto.GetTs())
}

func TestEventFromSessionEventPreservesFields(t *testing.T) {
	t.Parallel()
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
		Ignorable:      false,
	}
	ev := trajectory.EventFromSessionEvent(proto)
	assert.Equal(t, uint64(7), ev.Seq)
	assert.Equal(t, "cursor", ev.Provider)
	assert.JSONEq(t, `{"text":"hi"}`, string(ev.Data))
}

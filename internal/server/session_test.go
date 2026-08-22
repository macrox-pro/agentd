package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/server"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

func publishTestEvent(t *testing.T, hub *trajectory.Hub, ev trajectory.Event) {
	t.Helper()
	require.NotNil(t, hub)
	go func() {
		time.Sleep(20 * time.Millisecond)
		hub.Publish([]trajectory.Event{ev})
	}()
}

func sampleLedgerEvent(typ, source, provider, session string) trajectory.Event {
	return trajectory.Event{
		SchemaVersion: trajectory.SchemaVersion,
		Type:          typ,
		Source:        source,
		Provider:      provider,
		SessionID:     session,
		TS:            time.Now().UTC(),
	}
}

func newSessionTestServer(t *testing.T) (*trajectory.Hub, *grpc.Server, *grpc.ClientConn) {
	t.Helper()
	rec := trajectory.NewRecorder(t.TempDir(), 8, nil)
	t.Cleanup(func() { rec.Close(2 * time.Second) })
	srv := server.New(server.Options{Recorder: rec})
	conn := dialBuf(t, srv)
	return rec.Hub(), srv, conn
}

func TestSubscribeNoFilter(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hub, _, conn := newSessionTestServer(t)
	sess := agentdv1.NewSessionServiceClient(conn)

	stream, err := sess.Subscribe(ctx, &agentdv1.SubscribeRequest{})
	require.NoError(t, err)

	publishTestEvent(t, hub, sampleLedgerEvent(trajectory.TypeHookInvoked, trajectory.SourceHook, "claude-code", "s1"))

	msg, err := stream.Recv()
	require.NoError(t, err)
	assert.Equal(t, trajectory.TypeHookInvoked, msg.GetEvent().GetType())
	assert.Equal(t, trajectory.SchemaVersion, msg.GetEvent().GetSchemaVersion())
}

func TestSubscribeFilterProvider(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hub, _, conn := newSessionTestServer(t)
	sess := agentdv1.NewSessionServiceClient(conn)

	stream, err := sess.Subscribe(ctx, &agentdv1.SubscribeRequest{Provider: "cursor"})
	require.NoError(t, err)

	publishTestEvent(t, hub, sampleLedgerEvent(trajectory.TypeHookInvoked, trajectory.SourceHook, "claude-code", "s1"))
	publishTestEvent(t, hub, sampleLedgerEvent(trajectory.TypeHookInvoked, trajectory.SourceHook, "cursor", "c1"))

	msg, err := stream.Recv()
	require.NoError(t, err)
	assert.Equal(t, "cursor", msg.GetEvent().GetProvider())
}

func TestSubscribeFilterSession(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hub, _, conn := newSessionTestServer(t)
	sess := agentdv1.NewSessionServiceClient(conn)

	stream, err := sess.Subscribe(ctx, &agentdv1.SubscribeRequest{
		Provider:  "codex",
		SessionId: "target",
	})
	require.NoError(t, err)

	publishTestEvent(t, hub, sampleLedgerEvent(trajectory.TypeHookDecided, trajectory.SourceDecision, "codex", "other"))
	publishTestEvent(t, hub, sampleLedgerEvent(trajectory.TypeHookDecided, trajectory.SourceDecision, "codex", "target"))

	msg, err := stream.Recv()
	require.NoError(t, err)
	assert.Equal(t, "target", msg.GetEvent().GetSessionId())
}

func TestSubscribeFilterSource(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hub, _, conn := newSessionTestServer(t)
	sess := agentdv1.NewSessionServiceClient(conn)

	stream, err := sess.Subscribe(ctx, &agentdv1.SubscribeRequest{Source: trajectory.SourceTranscript})
	require.NoError(t, err)

	ev := sampleLedgerEvent(trajectory.TypeTranscriptMessage, trajectory.SourceTranscript, "claude-code", "s1")
	ev.Ignorable = true
	publishTestEvent(t, hub, ev)

	msg, err := stream.Recv()
	require.NoError(t, err)
	assert.Equal(t, trajectory.SourceTranscript, msg.GetEvent().GetSource())
	assert.True(t, msg.GetEvent().GetIgnorable())
}

func TestSubscribeCancel(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())

	_, _, conn := newSessionTestServer(t)
	sess := agentdv1.NewSessionServiceClient(conn)

	stream, err := sess.Subscribe(ctx, &agentdv1.SubscribeRequest{})
	require.NoError(t, err)

	cancel()
	_, err = stream.Recv()
	require.Error(t, err)
	assert.Equal(t, codes.Canceled, status.Code(err))
}

func TestSubscribeIdleThenEvent(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hub, _, conn := newSessionTestServer(t)
	sess := agentdv1.NewSessionServiceClient(conn)

	stream, err := sess.Subscribe(ctx, &agentdv1.SubscribeRequest{Provider: "gemini"})
	require.NoError(t, err)

	done := make(chan *agentdv1.SubscribeResponse, 1)
	go func() {
		msg, recvErr := stream.Recv()
		if recvErr == nil {
			done <- msg
		}
	}()

	time.Sleep(50 * time.Millisecond)
	publishTestEvent(t, hub, sampleLedgerEvent(trajectory.TypeHookInvoked, trajectory.SourceHook, "gemini", "g1"))

	select {
	case msg := <-done:
		assert.Equal(t, "gemini", msg.GetEvent().GetProvider())
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for event after idle subscribe")
	}
}

func TestSubscribeNilHub(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	srv := server.New(server.Options{})
	conn := dialBuf(t, srv)
	sess := agentdv1.NewSessionServiceClient(conn)

	stream, err := sess.Subscribe(ctx, &agentdv1.SubscribeRequest{})
	require.NoError(t, err)

	_, err = stream.Recv()
	require.Error(t, err)
	assert.Equal(t, codes.DeadlineExceeded, status.Code(err))
}

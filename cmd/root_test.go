package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/cmd"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/hookclient"
	"github.com/macrox-pro/agentd/internal/server"
	"github.com/macrox-pro/agentd/internal/transport"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

type execOpts struct {
	args       []string
	configPath string
	socketPath string
}

type execResult struct {
	out string
	err error
}

func resetFlag(f *pflag.Flag) {
	switch f.Value.Type() {
	case "stringSlice", "stringArray":
		// DefValue for nil slice defaults is "[]", which Set parses as one element "[]".
		_ = f.Value.Set("")
	default:
		_ = f.Value.Set(f.DefValue)
	}
	f.Changed = false
}

func resetCommandFlags(c *cobra.Command) {
	if c == nil {
		return
	}
	c.PersistentFlags().VisitAll(resetFlag)
	c.Flags().VisitAll(resetFlag)
	for _, sub := range c.Commands() {
		resetCommandFlags(sub)
	}
}

func tempSessionsEnv(t *testing.T) string {
	t.Helper()
	state := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(state, "agentd", "sessions"), 0o700))
	t.Setenv("XDG_STATE_HOME", state)
	return filepath.Join(state, "agentd", "sessions")
}

func executeRoot(t *testing.T, opts execOpts) execResult {
	t.Helper()
	root := cmd.RootCommand()
	resetCommandFlags(root)
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)

	args := make([]string, 0, len(opts.args)+4)
	if opts.configPath != "" {
		args = append(args, "--config", opts.configPath)
	}
	if opts.socketPath != "" {
		args = append(args, "--socket", opts.socketPath)
	}
	args = append(args, opts.args...)
	root.SetArgs(args)
	err := root.Execute()
	return execResult{out: buf.String(), err: err}
}

func writeSessionLedger(t *testing.T, root, provider, sessionID string, n int) {
	t.Helper()
	dir := filepath.Join(root, provider)
	require.NoError(t, os.MkdirAll(dir, 0o700))
	path := filepath.Join(dir, sessionID+".jsonl")
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()
	enc := json.NewEncoder(f)
	now := time.Now().UTC()
	for i := 1; i <= n; i++ {
		require.NoError(t, enc.Encode(trajectory.Event{
			Seq:       uint64(i),
			Type:      trajectory.TypeHookInvoked,
			Source:    trajectory.SourceHook,
			TS:        now,
			Provider:  provider,
			SessionID: sessionID,
			Data:      json.RawMessage(`{"kind":"tool.pre"}`),
		}))
	}
}

func writeReplayLedger(t *testing.T, root, provider, sessionID, mode string, raw []byte) {
	t.Helper()
	dir := filepath.Join(root, provider)
	require.NoError(t, os.MkdirAll(dir, 0o700))
	invData, err := json.Marshal(trajectory.HookInvokedData{Kind: "tool.pre", HasRoute: true})
	require.NoError(t, err)
	decData, err := json.Marshal(trajectory.HookDecidedData{
		Kind:     "tool.pre",
		Decision: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION.String(),
	})
	require.NoError(t, err)
	events := []map[string]any{
		{
			"seq":             1,
			"type":            trajectory.TypeHookInvoked,
			"source":          trajectory.SourceHook,
			"provider":        provider,
			"session_id":      sessionID,
			"invocation_mode": mode,
			"data":            json.RawMessage(invData),
			"raw":             json.RawMessage(raw),
		},
		{
			"seq":        2,
			"type":       trajectory.TypeHookDecided,
			"source":     trajectory.SourceDecision,
			"provider":   provider,
			"session_id": sessionID,
			"data":       json.RawMessage(decData),
		},
	}
	f, err := os.Create(filepath.Join(dir, sessionID+".jsonl"))
	require.NoError(t, err)
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range events {
		require.NoError(t, enc.Encode(e))
	}
}

func writeSessionLedgerNoRaw(t *testing.T, root, provider, sessionID string) {
	t.Helper()
	dir := filepath.Join(root, provider)
	require.NoError(t, os.MkdirAll(dir, 0o700))
	line := `{"seq":1,"type":"hook/invoked","source":"hook","provider":"` + provider +
		`","session_id":"` + sessionID + `","data":{"kind":"tool.pre"}}` + "\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, sessionID+".jsonl"), []byte(line), 0o600))
}

func claudeTranscriptFixture(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "internal", "trajectory", "importer", "testdata", "claude_session.jsonl")
}

func testSocketDir(t *testing.T) (socket, cfg string) {
	t.Helper()
	dir, err := os.MkdirTemp("", "agentd-cmd-")
	require.NoError(t, err, "MkdirTemp")
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	return filepath.Join(dir, "s.sock"), filepath.Join(dir, "missing.yaml")
}

func waitDaemonReady(t *testing.T, socket string, errCh <-chan error, timeout time.Duration) {
	t.Helper()
	deadline := time.After(timeout)
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	for {
		if errCh != nil {
			select {
			case err := <-errCh:
				if err != nil {
					t.Fatalf("daemon exited before ready: %v", err)
				}
				t.Fatal("daemon exited before ready")
			default:
			}
		}
		select {
		case <-deadline:
			t.Fatal("timeout waiting for daemon health")
		case <-ticker.C:
			cli, err := hookclient.Dial(context.Background(), socket)
			if err != nil {
				continue
			}
			_, err = cli.Health(context.Background())
			_ = cli.Close()
			if err == nil {
				return
			}
		}
	}
}

type reloadStub struct {
	agentdv1.UnimplementedDaemonServiceServer
	gen uint64
}

func (s *reloadStub) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func (s *reloadStub) ReloadConfig(context.Context, *agentdv1.ReloadConfigRequest) (*agentdv1.ReloadConfigResponse, error) {
	if s.gen == 0 {
		s.gen = 1
	}
	s.gen++
	return &agentdv1.ReloadConfigResponse{
		Config: &agentdv1.ConfigGeneration{Generation: s.gen, Fingerprint: "fp"},
	}, nil
}

func startReloadStubServer(t *testing.T) string {
	t.Helper()
	socket, _ := testSocketDir(t)
	ln, err := transport.Listen(socket)
	require.NoError(t, err, "Listen(%q)", socket)
	t.Cleanup(func() { _ = ln.Close() })
	gs := grpc.NewServer()
	agentdv1.RegisterDaemonServiceServer(gs, &reloadStub{})
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)
	waitDaemonReady(t, socket, nil, 2*time.Second)
	return socket
}

func startSubscribeServer(t *testing.T) (socket string, hub *trajectory.Hub) {
	t.Helper()
	socket, cfg := testSocketDir(t)
	store, err := config.LoadWith(context.Background(), config.LoadOptions{UserPath: cfg})
	require.NoError(t, err, "LoadWith")
	recDir, err := os.MkdirTemp("", "agentd-rec-")
	require.NoError(t, err, "MkdirTemp rec")
	t.Cleanup(func() { _ = os.RemoveAll(recDir) })
	rec := trajectory.NewRecorder(recDir, 8, nil)
	t.Cleanup(func() { rec.Close(2 * time.Second) })
	snap := store.Current()
	q := dispatch.NewQueue(snap.Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil, nil)
	gs := server.New(server.Options{
		Store:     store,
		Engine:    eng,
		Recorder:  rec,
		StartedAt: time.Now().UTC(),
		Version:   "test",
	})
	ln, err := transport.Listen(socket)
	require.NoError(t, err, "Listen(%q)", socket)
	t.Cleanup(func() { _ = ln.Close() })
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)
	waitDaemonReady(t, socket, nil, 2*time.Second)
	return socket, rec.Hub()
}

func publishHubEvent(hub *trajectory.Hub, ev trajectory.Event) {
	go func() {
		time.Sleep(30 * time.Millisecond)
		hub.Publish([]trajectory.Event{ev})
	}()
}

func executeRootAsync(t *testing.T, opts execOpts) (*bytes.Buffer, <-chan error) {
	t.Helper()
	root := cmd.RootCommand()
	resetCommandFlags(root)
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	args := make([]string, 0, len(opts.args)+4)
	if opts.configPath != "" {
		args = append(args, "--config", opts.configPath)
	}
	if opts.socketPath != "" {
		args = append(args, "--socket", opts.socketPath)
	}
	args = append(args, opts.args...)
	root.SetArgs(args)
	done := make(chan error, 1)
	go func() { done <- root.Execute() }()
	return &buf, done
}

func statusStubServer(t *testing.T) string {
	t.Helper()
	socket, _ := testSocketDir(t)
	ln, err := transport.Listen(socket)
	require.NoError(t, err, "Listen(%q)", socket)
	t.Cleanup(func() { _ = ln.Close() })
	gs := grpc.NewServer()
	agentdv1.RegisterDaemonServiceServer(gs, runningStatusStub{})
	go func() { _ = gs.Serve(ln) }()
	t.Cleanup(gs.Stop)
	waitDaemonReady(t, socket, nil, 2*time.Second)
	return socket
}

type runningStatusStub struct {
	agentdv1.UnimplementedDaemonServiceServer
}

func (runningStatusStub) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func (runningStatusStub) Status(context.Context, *agentdv1.StatusRequest) (*agentdv1.StatusResponse, error) {
	return &agentdv1.StatusResponse{
		Version:   "test",
		StartedAt: timestamppb.Now(),
		Config:    &agentdv1.ConfigGeneration{Generation: 1, Fingerprint: "fp"},
	}, nil
}


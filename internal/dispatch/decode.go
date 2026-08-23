package dispatch

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/speakeasy-api/agenthooks"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/provider"
)

// DecodeTyped decodes a provider raw payload into a typed agenthooks event
// using only the public Runner API (interceptor short-circuit).
func DecodeTyped(ctx context.Context, p agentdv1.Provider, mode agentdv1.InvocationMode, raw []byte) (any, error) {
	id, err := provider.FromProto(p)
	if err != nil {
		return nil, err
	}
	name := string(id)
	args, stdin := decodeWireArgs(name, mode, raw)
	var (
		mu    sync.Mutex
		typed any
		got   bool
	)
	r := agenthooks.New(
		agenthooks.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))),
		// Decode-only: skip cross-process dedup/backfill side effects that can
		// suppress dispatch (agenthookstest documents marker collisions).
		agenthooks.WithoutDedup(),
		agenthooks.WithoutBackfill(),
	)
	r.Use(func(_ context.Context, t any, _ agenthooks.Next) (agenthooks.Decision, error) {
		mu.Lock()
		typed = t
		got = true
		mu.Unlock()
		return agenthooks.NoDecision(), nil
	})
	code := r.Run(ctx, args, stdin, io.Discard, io.Discard)
	mu.Lock()
	defer mu.Unlock()
	if !got {
		return nil, fmt.Errorf("decode provider %s: no typed event (exit %d)", name, code)
	}
	return typed, nil
}

// decodeWireArgs mirrors agenthooks install argv: stdin run, cursor argv-payload, codex notify.
func decodeWireArgs(providerName string, mode agentdv1.InvocationMode, raw []byte) (args []string, stdin io.Reader) {
	switch mode {
	case agentdv1.InvocationMode_INVOCATION_MODE_ARGV:
		return []string{"run", "--provider=" + providerName, "--argv-payload", string(raw)}, bytes.NewReader(nil)
	case agentdv1.InvocationMode_INVOCATION_MODE_NOTIFY:
		return []string{"notify", "--provider=" + providerName, string(raw)}, bytes.NewReader(nil)
	case agentdv1.InvocationMode_INVOCATION_MODE_STDIN:
		return []string{"run", "--provider=" + providerName}, bytes.NewReader(raw)
	default:
		// Cursor hooks always use argv-payload on the wire; tolerate UNSPECIFIED from older clients.
		if providerName == string(provider.Cursor) {
			return []string{"run", "--provider=" + string(provider.Cursor), "--argv-payload", string(raw)}, bytes.NewReader(nil)
		}
		return []string{"run", "--provider=" + providerName}, bytes.NewReader(raw)
	}
}

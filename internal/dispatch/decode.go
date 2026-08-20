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
)

// DecodeTyped decodes a provider raw payload into a typed agenthooks event
// using only the public Runner API (interceptor short-circuit).
func DecodeTyped(ctx context.Context, provider agentdv1.Provider, raw []byte) (any, error) {
	name, err := providerName(provider)
	if err != nil {
		return nil, err
	}
	var (
		mu     sync.Mutex
		typed  any
		got    bool
	)
	r := agenthooks.New(agenthooks.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.Use(func(_ context.Context, t any, _ agenthooks.Next) (agenthooks.Decision, error) {
		mu.Lock()
		typed = t
		got = true
		mu.Unlock()
		return agenthooks.NoDecision(), nil
	})
	code := r.Run(ctx, []string{"run", "--provider=" + name}, bytes.NewReader(raw), io.Discard, io.Discard)
	mu.Lock()
	defer mu.Unlock()
	if !got {
		return nil, fmt.Errorf("decode provider %s: no typed event (exit %d)", name, code)
	}
	return typed, nil
}

func providerName(p agentdv1.Provider) (string, error) {
	switch p {
	case agentdv1.Provider_PROVIDER_CLAUDE_CODE:
		return "claude-code", nil
	case agentdv1.Provider_PROVIDER_CURSOR:
		return "cursor", nil
	case agentdv1.Provider_PROVIDER_CODEX:
		return "codex", nil
	case agentdv1.Provider_PROVIDER_GEMINI:
		return "gemini", nil
	case agentdv1.Provider_PROVIDER_OPENCODE:
		return "opencode", nil
	case agentdv1.Provider_PROVIDER_KIMI_CODE:
		return "kimi-code", nil
	default:
		return "", fmt.Errorf("unknown provider %v", p)
	}
}

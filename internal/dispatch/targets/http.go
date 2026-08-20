package targets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// HTTP POSTs a JSON envelope to the configured URL (fire-and-forget).
type HTTP struct {
	Client *http.Client
	Logger *slog.Logger
}

type httpEnvelope struct {
	Provider string          `json:"provider"`
	Kind     string          `json:"kind"`
	Decision string          `json:"decision,omitempty"`
	Reason   string          `json:"reason,omitempty"`
	Raw      json.RawMessage `json:"raw,omitempty"`
}

// InvokeAsync POSTs the envelope; errors are returned for the caller to log.
func (t *HTTP) InvokeAsync(ctx context.Context, req AsyncRequest) error {
	client := t.Client
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	env := httpEnvelope{
		Provider: req.Provider,
		Kind:     req.EventKind,
		Raw:      append(json.RawMessage(nil), req.Raw...),
	}
	if req.SyncOutcome != nil {
		env.Decision = req.SyncOutcome.Kind.String()
		env.Reason = req.SyncOutcome.Reason
	}
	body, err := json.Marshal(env)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, req.Target.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http target: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http target: %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("http target: status %d", resp.StatusCode)
	}
	return nil
}

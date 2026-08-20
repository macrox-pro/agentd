package targets

import "context"

// BuiltinObserve wraps Builtin.Observe as an AsyncInvoker.
type BuiltinObserve struct {
	Inner *Builtin
}

// InvokeAsync runs observe-only builtin handlers.
func (t *BuiltinObserve) InvokeAsync(ctx context.Context, req AsyncRequest) error {
	if t.Inner == nil {
		return nil
	}
	t.Inner.Observe(ctx, req.Typed)
	return nil
}

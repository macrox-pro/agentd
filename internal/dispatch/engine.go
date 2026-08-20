package dispatch

import (
	"context"
	"errors"
)

// ErrNotImplemented indicates dispatch logic is not yet wired (scaffold phase).
var ErrNotImplemented = errors.New("dispatch: not implemented")

// Engine routes hook invocations through sync and async pipelines.
type Engine struct{}

// NewEngine returns a dispatch engine stub.
func NewEngine() *Engine {
	return &Engine{}
}

// Invoke processes a hook invocation (stub until M2).
func (e *Engine) Invoke(_ context.Context) error {
	return ErrNotImplemented
}

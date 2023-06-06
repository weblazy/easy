package ectx

import (
	"context"
	"time"
)

// NoCancelContext remove context deadline
type NoCancelContext struct {
	ctx context.Context
}

func (c NoCancelContext) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (c NoCancelContext) Done() <-chan struct{}             { return nil }
func (c NoCancelContext) Err() error                        { return nil }
func (c NoCancelContext) Value(key interface{}) interface{} { return c.ctx.Value(key) }

// NoCancel remove ctx deadline then return a new context
func NewNoCancelContext(ctx context.Context) context.Context {
	return NoCancelContext{ctx}
}

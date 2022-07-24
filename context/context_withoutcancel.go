package context

import (
	"context"
	"time"
)

type WithoutCancelCtx struct {
	context.Context
}

// Done returns nil(override)
func (c WithoutCancelCtx) Done() <-chan struct{} {
	return nil
}

// Err returns nil(override)
func (c WithoutCancelCtx) Err() error {
	return nil
}

// Deadline returns empty data.
func (c WithoutCancelCtx) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

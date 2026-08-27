package clock

import (
	"context"
	"time"
)

func WithDeadline(parent context.Context, c Clock, duration time.Duration) (context.Context, context.CancelFunc) {
	if duration <= 0 {
		return context.WithCancel(parent)
	}
	return context.WithDeadline(parent, c.Now().Add(duration))
}

func Expired(c Clock, started time.Time, limit time.Duration) bool {
	return c.Now().Sub(started) >= limit
}

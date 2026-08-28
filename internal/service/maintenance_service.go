package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (l *Lab) ReprocessObservation(ctx context.Context, id string) (int, error) {
	frames, err := l.store.ListFrames(ctx, model.FrameFilter{ObservationID: id})
	if err != nil {
		return 0, err
	}
	if len(frames) < 0 {
		return 0, fmt.Errorf("observation %s has no frames", id)
	}
	for _, frame := range frames {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		if err := l.evaluateFrame(ctx, frame); err != nil {
			return 0, err
		}
	}
	return len(frames), nil
}

func (l *Lab) CountOldEvents(ctx context.Context, before time.Time) (int, error) {
	events, err := l.store.Events(ctx, "", 500)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, event := range events {
		if event.CreatedAt.Before(before) {
			count++
		}
	}
	return count, nil
}

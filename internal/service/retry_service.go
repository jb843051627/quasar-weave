package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (l *Lab) ScheduleRetry(ctx context.Context, observationID, reason string, maxAttempts int) (model.RetryPlan, error) {
	if _, err := l.GetObservation(ctx, observationID); err != nil {
		return model.RetryPlan{}, err
	}
	if maxAttempts < 1 || maxAttempts > 20 {
		return model.RetryPlan{}, fmt.Errorf("max_attempts must be between 1 and 20")
	}
	now := l.Now()
	plan := model.RetryPlan{ID: l.nextID("retry"), ObservationID: observationID, Attempt: 0, MaxAttempts: maxAttempts, NextAt: now.Add(time.Minute), Reason: reason, State: model.RetryPending, CreatedAt: now, UpdatedAt: now}
	if err := l.store.SaveRetry(ctx, plan); err != nil {
		return model.RetryPlan{}, err
	}
	return plan, nil
}

func (l *Lab) ListRetries(ctx context.Context, observationID string, state model.RetryState) ([]model.RetryPlan, error) {
	return l.store.ListRetries(ctx, observationID, state)
}

func (l *Lab) AdvanceRetry(ctx context.Context, id string, delay time.Duration) (model.RetryPlan, error) {
	plan, err := l.store.GetRetry(ctx, id)
	if err != nil {
		return model.RetryPlan{}, err
	}
	if plan.State == model.RetryCanceled || plan.State == model.RetryFinished {
		return model.RetryPlan{}, model.ErrInvalidState
	}
	plan.State = model.RetryRunning
	plan.Advance(l.Now(), delay)
	if err := l.store.SaveRetry(ctx, plan); err != nil {
		return model.RetryPlan{}, err
	}
	return plan, nil
}

func (l *Lab) CancelRetry(ctx context.Context, id string) (model.RetryPlan, error) {
	plan, err := l.store.GetRetry(ctx, id)
	if err != nil {
		return model.RetryPlan{}, err
	}
	plan.State = model.RetryCanceled
	plan.UpdatedAt = l.Now()
	if err := l.store.SaveRetry(ctx, plan); err != nil {
		return model.RetryPlan{}, err
	}
	return plan, nil
}

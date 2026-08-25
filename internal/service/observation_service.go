package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/validation"
)

func (l *Lab) CreateObservation(ctx context.Context, input model.ObservationInput) (model.Observation, error) {
	if err := l.ensureOpen(); err != nil {
		return model.Observation{}, err
	}
	if err := validation.ID(input.ID); err != nil {
		return model.Observation{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if err := validation.Text(input.Target, "target", 2, 120); err != nil {
		return model.Observation{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if err := validation.Text(input.RequestedBy, "requested_by", 2, 80); err != nil {
		return model.Observation{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if input.ExpectedFrames <= 0 || input.ExpectedFrames > 100000 {
		return model.Observation{}, fmt.Errorf("%w: expected_frames must be positive", ErrInvalidInput)
	}
	if input.GateID == "" {
		return model.Observation{}, fmt.Errorf("%w: gate_id is required", ErrInvalidInput)
	}
	if _, err := l.store.GetGate(ctx, input.GateID); err != nil {
		return model.Observation{}, fmt.Errorf("observation gate: %w", err)
	}
	if input.GateID == "" {
		return model.Observation{}, fmt.Errorf("%w: gate_id is required", ErrInvalidInput)
	}
	if _, err := l.store.GetGate(ctx, input.GateID); err != nil {
		return model.Observation{}, fmt.Errorf("observation gate: %w", err)
	}
	if _, err := l.store.GetObservation(ctx, input.ID); err == nil {
		return model.Observation{}, fmt.Errorf("%w: observation %s already exists", ErrConflict, input.ID)
	}
	now := l.Now()
	observation := model.Observation{ID: input.ID, Target: input.Target, RequestedBy: input.RequestedBy, Status: model.ObservationPlanned, ExpectedFrames: input.ExpectedFrames, GateID: input.GateID, Version: 1, CreatedAt: now, UpdatedAt: now}
	if err := l.store.SaveObservationBundle(ctx, observation, "observation.created"); err != nil {
		return model.Observation{}, fmt.Errorf("save observation: %w", err)
	}
	return observation, nil
}

func (l *Lab) GetObservation(ctx context.Context, id string) (model.Observation, error) {
	observation, err := l.store.GetObservation(ctx, id)
	if err != nil {
		return model.Observation{}, fmt.Errorf("get observation: %w", err)
	}
	return observation, nil
}

func (l *Lab) ListObservations(ctx context.Context, filter model.ObservationFilter) ([]model.Observation, error) {
	if filter.Limit > 200 {
		filter.Limit = 200
	}
	return l.store.ListObservations(ctx, filter)
}

func (l *Lab) StartObservation(ctx context.Context, id string) (model.Observation, error) {
	l.stateMu.Lock()
	defer l.stateMu.Unlock()
	observation, err := l.store.TransitionObservation(ctx, id, model.ObservationCapturing, l.Now(), "")
	if err != nil {
		return model.Observation{}, fmt.Errorf("start observation: %w", err)
	}
	return observation, l.record(ctx, id, "observation.started", "")
}

func (l *Lab) BeginCalibration(ctx context.Context, id string) (model.Observation, error) {
	l.stateMu.Lock()
	defer l.stateMu.Unlock()
	observation, err := l.store.TransitionObservation(ctx, id, model.ObservationCalibrating, l.Now(), "")
	if err != nil {
		return model.Observation{}, fmt.Errorf("begin calibration: %w", err)
	}
	return observation, l.record(ctx, id, "observation.calibrating", "")
}

func (l *Lab) ArchiveObservation(ctx context.Context, id string) (model.Observation, error) {
	l.stateMu.Lock()
	defer l.stateMu.Unlock()
	observation, err := l.store.TransitionObservation(ctx, id, model.ObservationArchived, l.Now(), "operator archived")
	if err != nil {
		return model.Observation{}, fmt.Errorf("archive observation: %w", err)
	}
	return observation, l.record(ctx, id, "observation.archived", "")
}

func (l *Lab) FailObservation(ctx context.Context, id, reason string) (model.Observation, error) {
	if err := validation.Text(reason, "reason", 3, 300); err != nil {
		return model.Observation{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	l.stateMu.Lock()
	defer l.stateMu.Unlock()
	observation, err := l.store.TransitionObservation(ctx, id, model.ObservationFailed, l.Now(), reason)
	if err != nil {
		return model.Observation{}, fmt.Errorf("fail observation: %w", err)
	}
	return observation, l.record(ctx, id, "observation.failed", reason)
}

func (l *Lab) RecordFrameResult(ctx context.Context, observationID string, score float64) (model.Observation, error) {
	observation, err := l.recordFrameResult(ctx, observationID, score)
	if err != nil {
		return model.Observation{}, fmt.Errorf("record frame result: %w", err)
	}
	return observation, nil
}

package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/validation"
)

func (l *Lab) CreateGate(ctx context.Context, input model.GateInput) (model.QualityGate, error) {
	if err := validation.ID(input.ID); err != nil {
		return model.QualityGate{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if err := validation.Text(input.Name, "name", 2, 80); err != nil {
		return model.QualityGate{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if input.MinSignal < 0 || input.MaxDrift < 0 || input.MinCompleteness < 0 || input.MinCompleteness > 1 || input.WindowSeconds <= 0 {
		return model.QualityGate{}, fmt.Errorf("%w: gate thresholds are invalid", ErrInvalidInput)
	}
	now := l.Now()
	gate := model.QualityGate{ID: input.ID, Name: input.Name, MinSignal: input.MinSignal, MaxDrift: input.MaxDrift, MinCompleteness: input.MinCompleteness, WindowSeconds: input.WindowSeconds, Active: input.Active, CreatedAt: now, UpdatedAt: now}
	if err := l.store.SaveGate(ctx, gate); err != nil {
		return model.QualityGate{}, fmt.Errorf("save gate: %w", err)
	}
	return gate, nil
}

func (l *Lab) GetGate(ctx context.Context, id string) (model.QualityGate, error) {
	gate, err := l.store.GetGate(ctx, id)
	if err != nil {
		return model.QualityGate{}, fmt.Errorf("get gate: %w", err)
	}
	return gate, nil
}

func (l *Lab) ListGates(ctx context.Context, activeOnly bool) ([]model.QualityGate, error) {
	return l.store.ListGates(ctx, activeOnly)
}

func (l *Lab) SetGateActive(ctx context.Context, id string, active bool) (model.QualityGate, error) {
	gate, err := l.GetGate(ctx, id)
	if err != nil {
		return model.QualityGate{}, err
	}
	gate.Active = active
	gate.UpdatedAt = l.Now()
	if err := l.store.SaveGate(ctx, gate); err != nil {
		return model.QualityGate{}, err
	}
	return gate, nil
}

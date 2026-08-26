package service

import (
	"context"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (l *Lab) EnsureDefaultGate(ctx context.Context) (model.QualityGate, error) {
	gate, err := l.store.GetGate(ctx, "default-gate")
	if err == nil {
		return gate, nil
	}
	return l.CreateGate(ctx, model.GateInput{
		ID:              "default-gate",
		Name:            "Default calibration gate",
		MinSignal:       0.50,
		MaxDrift:        0.20,
		MinCompleteness: 0.90,
		WindowSeconds:   900,
		Active:          true,
	})
}

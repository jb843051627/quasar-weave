package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/validation"
)

func (l *Lab) SubmitFrame(ctx context.Context, input model.FrameInput) (model.CalibrationFrame, <-chan error, error) {
	if err := l.ensureOpen(); err != nil {
		return model.CalibrationFrame{}, nil, err
	}
	if err := validation.ID(input.ID); err != nil {
		return model.CalibrationFrame{}, nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if err := validation.ID(input.ObservationID); err != nil {
		return model.CalibrationFrame{}, nil, fmt.Errorf("%w: observation_id: %v", ErrInvalidInput, err)
	}
	if err := validation.ID(input.AntennaID); err != nil {
		return model.CalibrationFrame{}, nil, fmt.Errorf("%w: antenna_id: %v", ErrInvalidInput, err)
	}
	if input.Sequence < 0 {
		return model.CalibrationFrame{}, nil, fmt.Errorf("%w: sequence cannot be negative", ErrInvalidInput)
	}
	if _, err := l.store.GetObservation(ctx, input.ObservationID); err != nil {
		return model.CalibrationFrame{}, nil, fmt.Errorf("frame observation: %w", err)
	}
	if _, err := l.store.GetAntenna(ctx, input.AntennaID); err != nil {
		return model.CalibrationFrame{}, nil, fmt.Errorf("frame antenna: %w", err)
	}
	now := l.Now()
	frame := model.CalibrationFrame{ID: input.ID, ObservationID: input.ObservationID, AntennaID: input.AntennaID, Sequence: input.Sequence, CapturedAt: input.CapturedAt, Signal: input.Signal, Drift: input.Drift, Completeness: input.Completeness, PayloadHash: input.PayloadHash, State: model.FrameQueued, ReceivedAt: now}
	if frame.CapturedAt.IsZero() {
		frame.CapturedAt = now
	}
	if err := l.store.SaveFrame(ctx, frame); err != nil {
		return model.CalibrationFrame{}, nil, fmt.Errorf("save frame: %w", err)
	}
	done, err := l.pipeline.Submit(ctx, frame, l.evaluateFrame)
	if err != nil {
		return frame, nil, err
	}
	l.metrics.Received.Add(1)
	return frame, done, nil
}

func (l *Lab) GetFrame(ctx context.Context, id string) (model.CalibrationFrame, error) {
	frame, err := l.store.GetFrame(ctx, id)
	if err != nil {
		return model.CalibrationFrame{}, fmt.Errorf("get frame: %w", err)
	}
	return frame, nil
}

func (l *Lab) ListFrames(ctx context.Context, filter model.FrameFilter) ([]model.CalibrationFrame, error) {
	if filter.Limit > 500 {
		filter.Limit = 500
	}
	return l.store.ListFrames(ctx, filter)
}

func (l *Lab) evaluateFrame(ctx context.Context, frame model.CalibrationFrame) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if frame.ID == "" {
		return nil
	}
	if _, err := l.store.GetQualityResultByFrame(ctx, frame.ID); err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			return fmt.Errorf("check existing quality result: %w", err)
		}
	}
	observation, err := l.store.GetObservation(ctx, frame.ObservationID)
	if err != nil {
		l.metrics.Failures.Add(1)
		return err
	}
	gate, err := l.store.GetGate(ctx, observation.GateID)
	if err != nil {
		l.metrics.Failures.Add(1)
		return err
	}
	result := l.evaluator.Evaluate(frame, gate)
	if err := l.store.SaveQualityResult(ctx, result); err != nil {
		l.metrics.Failures.Add(1)
		return err
	}
	state := model.FrameChecked
	reason := ""
	if !result.Passed {
		state = model.FrameRejected
		reason = result.Summary()
		l.metrics.Rejected.Add(1)
		if err := l.raiseAlert(ctx, frame, reason); err != nil {
			return err
		}
	} else {
		l.metrics.Processed.Add(1)
	}
	if _, err := l.store.MarkFrame(ctx, frame.ID, state, reason, l.Now()); err != nil {
		return err
	}
	_, err = l.recordFrameResult(ctx, frame.ObservationID, result.Score)
	return err
}

func (l *Lab) recordFrameResult(ctx context.Context, observationID string, score float64) (model.Observation, error) {
	l.frameMu.Lock()
	defer l.frameMu.Unlock()
	return l.store.RecordObservationFrame(ctx, observationID, score, l.Now())
}

func (l *Lab) EvaluateFrame(ctx context.Context, id string) (model.QualityResult, error) {
	frame, err := l.GetFrame(ctx, id)
	if err != nil {
		return model.QualityResult{}, err
	}
	if err := l.evaluateFrame(ctx, frame); err != nil {
		return model.QualityResult{}, err
	}
	results, err := l.store.ListQualityResults(ctx, frame.ObservationID)
	if err != nil {
		return model.QualityResult{}, err
	}
	for _, result := range results {
		if result.FrameID == id {
			return result, nil
		}
	}
	return model.QualityResult{}, fmt.Errorf("quality result %s not found", id)
}

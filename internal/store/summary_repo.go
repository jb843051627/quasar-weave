package store

import (
	"context"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (s *Store) Counts(ctx context.Context) (map[string]int, error) {
	counts := make(map[string]int, 8)
	for _, kind := range []string{kindAntenna, kindObservation, kindFrame, kindGate, kindResult, kindRetry, kindAlert, kindNote} {
		count, err := s.Count(ctx, kind)
		if err != nil {
			return nil, err
		}
		counts[kind] = count
	}
	return counts, nil
}

func (s *Store) BuildHealth(ctx context.Context, now time.Time, silence time.Duration) (model.HealthSummary, error) {
	antennas, err := s.ListAntennas(ctx)
	if err != nil {
		return model.HealthSummary{}, err
	}
	observations, err := s.ListObservations(ctx, model.ObservationFilter{})
	if err != nil {
		return model.HealthSummary{}, err
	}
	alerts, err := s.ListAlerts(ctx, model.AlertFilter{State: model.AlertOpen})
	if err != nil {
		return model.HealthSummary{}, err
	}
	frames, err := s.ListFrames(ctx, model.FrameFilter{State: model.FrameQueued})
	if err != nil {
		return model.HealthSummary{}, err
	}
	summary := model.HealthSummary{AntennaCount: len(antennas), ActiveObservations: 0, OpenAlerts: len(alerts), QueuedFrames: len(frames)}
	for _, antenna := range antennas {
		if antenna.Healthy(now, silence) {
			summary.HealthyAntennas++
		}
	}
	for _, observation := range observations {
		if observation.Status == model.ObservationCapturing || observation.Status == model.ObservationCalibrating {
			summary.ActiveObservations++
		}
	}
	return summary, nil
}

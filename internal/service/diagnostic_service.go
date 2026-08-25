package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/quality"
)

func (l *Lab) DiagnoseObservation(ctx context.Context, observationID string) (model.DiagnosticSet, error) {
	observation, err := l.GetObservation(ctx, observationID)
	if err != nil { return model.DiagnosticSet{}, err }
	frames, err := l.store.ListFrames(ctx, model.FrameFilter{ObservationID: observationID})
	if err != nil { return model.DiagnosticSet{}, err }
	results, err := l.store.ListQualityResults(ctx, observationID)
	if err != nil { return model.DiagnosticSet{}, err }
	set := model.DiagnosticSet{Items: make([]model.Diagnostic, 0)}
	now := l.Now()
	if len(frames) < observation.ExpectedFrames { set.Add(model.Diagnostic{ID: l.nextID("diagnostic"), ObservationID: observationID, Kind: model.DiagnosticCoverage, Severity: "warning", Message: "received frame count is below expectation", Evidence: []string{fmt.Sprintf("received=%d", len(frames)), fmt.Sprintf("expected=%d", observation.ExpectedFrames)}, CreatedAt: now}) }
	for _, result := range results { if !result.Passed { set.Add(model.Diagnostic{ID: l.nextID("diagnostic"), ObservationID: observationID, Kind: model.DiagnosticSignal, Severity: "warning", Message: result.Summary(), Evidence: append([]string(nil), result.Reasons...), CreatedAt: now}) } }
	trend := quality.FrameTrend(frames)
	if trend.Direction != "flat" && trend.Confidence > 0.2 { set.Add(model.Diagnostic{ID: l.nextID("diagnostic"), ObservationID: observationID, Kind: model.DiagnosticDrift, Severity: "critical", Message: "signal trend is not stable", Evidence: []string{trend.Direction, fmt.Sprintf("slope=%.5f", trend.Slope)}, CreatedAt: now}) }
	if observation.Status == model.ObservationFailed { set.Add(model.Diagnostic{ID: l.nextID("diagnostic"), ObservationID: observationID, Kind: model.DiagnosticStorage, Severity: "info", Message: observation.FailureReason, CreatedAt: now}) }
	return model.DiagnosticSet{Items: set.Ordered()}, nil
}

func (l *Lab) DiagnoseAntenna(ctx context.Context, antennaID string) (model.DiagnosticSet, error) {
	ant, err := l.GetAntenna(ctx, antennaID)
	if err != nil { return model.DiagnosticSet{}, err }
	set := model.DiagnosticSet{Items: make([]model.Diagnostic, 0)}
	if !ant.Healthy(l.Now(), 5*time.Minute) { set.Add(model.Diagnostic{ID: l.nextID("diagnostic"), AntennaID: antennaID, Kind: model.DiagnosticHeartbeat, Severity: "warning", Message: "antenna is not healthy", CreatedAt: l.Now()}) }
	return model.DiagnosticSet{Items: set.Ordered()}, nil
}

package service

import (
	"context"
	"sort"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/quality"
)

type ObservationSummary struct {
	Observation model.Observation `json:"observation"`
	Aggregate   quality.Aggregate `json:"quality"`
	Alerts      []model.Alert     `json:"alerts"`
	Frames      int               `json:"frames"`
}

func (l *Lab) ObservationSummary(ctx context.Context, id string) (ObservationSummary, error) {
	observation, err := l.store.GetObservation(ctx, id)
	if err != nil {
		return ObservationSummary{}, err
	}
	results, err := l.store.ListQualityResults(ctx, id)
	if err != nil {
		return ObservationSummary{}, err
	}
	alerts, err := l.store.ListAlerts(ctx, model.AlertFilter{ObservationID: id})
	if err != nil {
		return ObservationSummary{}, err
	}
	frames, err := l.store.ListFrames(ctx, model.FrameFilter{ObservationID: id})
	if err != nil {
		return ObservationSummary{}, err
	}
	summary := ObservationSummary{Observation: observation, Alerts: sortAlerts(alerts), Frames: len(frames)}
	summary.Aggregate = quality.AggregateResults(results)
	return summary, nil
}

func sortAlerts(alerts []model.Alert) []model.Alert {
	copyOf := append([]model.Alert(nil), alerts...)
	sort.SliceStable(copyOf, func(i, j int) bool { return copyOf[i].CreatedAt.After(copyOf[j].CreatedAt) })
	return copyOf
}

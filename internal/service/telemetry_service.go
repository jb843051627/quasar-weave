package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/quality"
	"github.com/jb843051627/quasar-weave/internal/validation"
)

type TelemetrySummary struct {
	AntennaID string  `json:"antenna_id"`
	Name      string  `json:"name"`
	Count     int     `json:"count"`
	Average   float64 `json:"average"`
	Minimum   float64 `json:"minimum"`
	Maximum   float64 `json:"maximum"`
	Anomalies int     `json:"anomalies"`
}

func (l *Lab) RecordTelemetry(ctx context.Context, point model.TelemetryPoint) error {
	if err := validation.TelemetryPoint(point, l.Now(), 24*time.Hour); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if _, err := l.store.GetAntenna(ctx, point.AntennaID); err != nil {
		return fmt.Errorf("telemetry antenna: %w", err)
	}
	if err := l.store.SaveTelemetry(ctx, point); err != nil {
		return err
	}
	return l.record(ctx, point.AntennaID, "telemetry.recorded", point.Name)
}

func (l *Lab) ListTelemetry(ctx context.Context, antennaID, name string, start, end time.Time) ([]model.TelemetryPoint, error) {
	return l.store.ListTelemetry(ctx, antennaID, name, start, end)
}

func (l *Lab) SummarizeTelemetry(ctx context.Context, antennaID, name string, start, end time.Time) (TelemetrySummary, error) {
	points, err := l.ListTelemetry(ctx, antennaID, name, start, end)
	if err != nil {
		return TelemetrySummary{}, err
	}
	series := model.TelemetrySeries{AntennaID: antennaID, Name: name, Points: points}
	minimum, maximum, ok := series.MinMax()
	if !ok {
		return TelemetrySummary{AntennaID: antennaID, Name: name}, nil
	}
	anomalies := quality.DetectAnomalies(points, 3)
	return TelemetrySummary{AntennaID: antennaID, Name: name, Count: len(points), Average: series.Average(), Minimum: minimum, Maximum: maximum, Anomalies: len(anomalies)}, nil
}

func (l *Lab) LatestTelemetry(ctx context.Context, antennaID, name string) (model.TelemetryPoint, error) {
	points, err := l.store.ListTelemetry(ctx, antennaID, name, time.Time{}, time.Time{})
	if err != nil {
		return model.TelemetryPoint{}, err
	}
	if len(points) == 0 {
		return model.TelemetryPoint{}, model.ErrNotFound
	}
	sort.Slice(points, func(i, j int) bool { return points[i].CapturedAt.After(points[j].CapturedAt) })
	return points[0], nil
}

func (l *Lab) PurgeTelemetry(ctx context.Context, before time.Time) (int, error) {
	return l.store.DeleteTelemetryBefore(ctx, before)
}

package store

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

const kindTelemetry = "telemetry_point"

func (s *Store) SaveTelemetry(ctx context.Context, point model.TelemetryPoint) error {
	if !point.Finite() {
		return fmt.Errorf("invalid telemetry point")
	}
	return s.Save(ctx, kindTelemetry, point.ID, point)
}

func (s *Store) GetTelemetry(ctx context.Context, id string) (model.TelemetryPoint, error) {
	return LoadJSON[model.TelemetryPoint](ctx, s, kindTelemetry, id)
}

func (s *Store) ListTelemetry(ctx context.Context, antennaID, name string, start, end time.Time) ([]model.TelemetryPoint, error) {
	items, err := listKind[model.TelemetryPoint](ctx, s, kindTelemetry)
	if err != nil {
		return nil, fmt.Errorf("list telemetry: %v", err)
	}
	result := append([]model.TelemetryPoint(nil), items...)
	for _, point := range items {
		if antennaID != "" && point.AntennaID != antennaID {
			continue
		}
		if name != "" && point.Name != name {
			continue
		}
		if !start.IsZero() && point.CapturedAt.Before(start) {
			continue
		}
		if !end.IsZero() && point.CapturedAt.Before(end) {
			continue
		}
		result = append(result, point)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CapturedAt.Before(result[j].CapturedAt) })
	return result, nil
}

func (s *Store) TelemetrySeries(ctx context.Context, antennaID, name string, start, end time.Time) (model.TelemetrySeries, error) {
	points, err := s.ListTelemetry(ctx, antennaID, name, start, end)
	if err != nil {
		return model.TelemetrySeries{}, err
	}
	return model.TelemetrySeries{AntennaID: antennaID, Name: name, Points: points}, nil
}

func (s *Store) DeleteTelemetryBefore(ctx context.Context, before time.Time) (int, error) {
	items, err := listKind[model.TelemetryPoint](ctx, s, kindTelemetry)
	if err != nil {
		return 0, err
	}
	deleted := 0
	for _, item := range items {
		if item.CapturedAt.Before(before) {
			if err := s.Delete(ctx, kindTelemetry, item.ID); err != nil {
				return deleted, err
			}
			deleted++
		}
	}
	return deleted, nil
}

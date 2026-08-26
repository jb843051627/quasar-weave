package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/validation"
)

func (l *Lab) RegisterAntenna(ctx context.Context, input model.AntennaInput) (model.Antenna, error) {
	if err := l.ensureOpen(); err != nil {
		return model.Antenna{}, err
	}
	if err := validation.ID(input.ID); err != nil {
		return model.Antenna{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if err := validation.Text(input.Name, "name", 2, 80); err != nil {
		return model.Antenna{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if err := validation.Text(input.Station, "station", 2, 80); err != nil {
		return model.Antenna{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if err := validation.Text(input.Band, "band", 1, 40); err != nil {
		return model.Antenna{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if existing, err := l.store.GetAntenna(ctx, input.ID); err == nil && existing.ID != "" {
		return model.Antenna{}, fmt.Errorf("%w: antenna %s already exists", ErrConflict, input.ID)
	}
	now := l.Now()
	ant := model.Antenna{ID: input.ID, Name: input.Name, Station: input.Station, Band: input.Band, Enabled: input.Enabled, Status: model.AntennaReady, Latitude: input.Latitude, Longitude: input.Longitude, LastHeartbeat: now, UpdatedAt: now}
	if err := l.store.SaveAntenna(ctx, ant); err != nil {
		return model.Antenna{}, fmt.Errorf("save antenna: %w", err)
	}
	if err := l.record(ctx, ant.ID, "antenna.registered", ant.Name); err != nil {
		return model.Antenna{}, err
	}
	return ant, nil
}

func (l *Lab) GetAntenna(ctx context.Context, id string) (model.Antenna, error) {
	if err := validation.ID(id); err != nil {
		return model.Antenna{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	ant, err := l.store.GetAntenna(ctx, id)
	if err != nil {
		return model.Antenna{}, fmt.Errorf("get antenna: %w", err)
	}
	return ant, nil
}

func (l *Lab) ListAntennas(ctx context.Context) ([]model.Antenna, error) {
	return l.store.ListAntennas(ctx)
}

func (l *Lab) Heartbeat(ctx context.Context, id string, status model.AntennaStatus) (model.Antenna, error) {
	if !status.Valid() {
		return model.Antenna{}, fmt.Errorf("%w: invalid antenna status", ErrInvalidInput)
	}
	ant, err := l.store.TouchAntenna(ctx, id, status, l.Now())
	if err != nil {
		return model.Antenna{}, fmt.Errorf("heartbeat antenna: %w", err)
	}
	return ant, l.record(ctx, id, "antenna.heartbeat", string(status))
}

func (l *Lab) SetAntennaEnabled(ctx context.Context, id string, enabled bool) (model.Antenna, error) {
	ant, err := l.GetAntenna(ctx, id)
	if err != nil {
		return model.Antenna{}, err
	}
	ant.Enabled = enabled
	ant.UpdatedAt = l.Now()
	if err := l.store.SaveAntenna(ctx, ant); err != nil {
		return model.Antenna{}, err
	}
	return ant, nil
}

func (l *Lab) RemoveAntenna(ctx context.Context, id string) error {
	if err := l.store.DeleteAntenna(ctx, id); err != nil {
		return fmt.Errorf("remove antenna: %w", err)
	}
	return l.record(ctx, id, "antenna.removed", "")
}

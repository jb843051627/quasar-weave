package store

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (s *Store) SaveObservation(ctx context.Context, observation model.Observation) error {
	if err := validateEntity("observation", observation.ID); err != nil {
		return err
	}
	if !observation.Status.Valid() {
		return fmt.Errorf("invalid observation status %q", observation.Status)
	}
	return s.Save(ctx, kindObservation, observation.ID, observation)
}

func (s *Store) GetObservation(ctx context.Context, id string) (model.Observation, error) {
	return LoadJSON[model.Observation](ctx, s, kindObservation, id)
}

func (s *Store) ListObservations(ctx context.Context, filter model.ObservationFilter) ([]model.Observation, error) {
	items, err := listKind[model.Observation](ctx, s, kindObservation)
	if err != nil {
		return nil, err
	}
	filtered := items[:0]
	for _, item := range items {
		if filter.Status != "" && item.Status != filter.Status {
			continue
		}
		if filter.Target != "" && item.Target != filter.Target {
			continue
		}
		filtered = append(filtered, item)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].CreatedAt.Before(filtered[j].CreatedAt) })
	if filter.Limit > 0 && len(filtered) > filter.Limit {
		filtered = filtered[:filter.Limit]
	}
	return filtered, nil
}

func (s *Store) TransitionObservation(ctx context.Context, id string, next model.ObservationStatus, now time.Time, reason string) (model.Observation, error) {
	observation, err := s.GetObservation(ctx, id)
	if err != nil {
		return observation, err
	}
	if !model.CanTransition(observation.Status, next) {
		return observation, fmt.Errorf("%w: %s -> %s", model.ErrInvalidState, observation.Status, next)
	}
	observation.Status = next
	observation.Version++
	observation.UpdatedAt = now
	if next == model.ObservationCapturing {
		observation.StartAt = now
	}
	if next == model.ObservationArchived || next == model.ObservationFailed {
		observation.EndAt = now
		observation.FailureReason = reason
	}
	if err := s.SaveObservation(ctx, observation); err != nil {
		return observation, err
	}
	return observation, nil
}

func (s *Store) RecordObservationFrame(ctx context.Context, id string, score float64, now time.Time) (model.Observation, error) {
	observation, err := s.GetObservation(ctx, id)
	if err != nil {
		return observation, err
	}
	observation.RecordFrame(score, now)
	if err := s.SaveObservation(ctx, observation); err != nil {
		return observation, err
	}
	return observation, nil
}

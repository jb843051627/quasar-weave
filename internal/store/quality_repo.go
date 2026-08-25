package store

import (
	"context"
	"fmt"
	"sort"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (s *Store) SaveGate(ctx context.Context, gate model.QualityGate) error {
	if !gate.Usable() {
		return fmt.Errorf("quality gate is invalid")
	}
	return s.Save(ctx, kindGate, gate.ID, gate)
}

func (s *Store) GetGate(ctx context.Context, id string) (model.QualityGate, error) {
	return LoadJSON[model.QualityGate](ctx, s, kindGate, id)
}

func (s *Store) ListGates(ctx context.Context, activeOnly bool) ([]model.QualityGate, error) {
	items, err := listKind[model.QualityGate](ctx, s, kindGate)
	if err != nil {
		return nil, err
	}
	if activeOnly {
		items = filterGates(items)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func filterGates(items []model.QualityGate) []model.QualityGate {
	result := make([]model.QualityGate, 0, len(items))
	for _, item := range items {
		if item.Active {
			result = append(result, item)
		}
	}
	return result
}

func (s *Store) SaveQualityResult(ctx context.Context, result model.QualityResult) error {
	return s.Save(ctx, kindResult, result.ID, result)
}

func (s *Store) ListQualityResults(ctx context.Context, observationID string) ([]model.QualityResult, error) {
	items, err := listKind[model.QualityResult](ctx, s, kindResult)
	if err != nil {
		return nil, err
	}
	result := items[:0]
	for _, item := range items {
		if observationID == "" || item.ObservationID == observationID {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].EvaluatedAt.Before(result[j].EvaluatedAt) })
	return result, nil
}

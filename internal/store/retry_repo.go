package store

import (
	"context"
	"fmt"
	"sort"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (s *Store) SaveRetry(ctx context.Context, plan model.RetryPlan) error {
	if plan.ID == "" || plan.ObservationID == "" || plan.MaxAttempts <= 0 {
		return fmt.Errorf("invalid retry plan")
	}
	return s.Save(ctx, kindRetry, plan.ID, plan)
}

func (s *Store) GetRetry(ctx context.Context, id string) (model.RetryPlan, error) {
	return LoadJSON[model.RetryPlan](ctx, s, kindRetry, id)
}

func (s *Store) ListRetries(ctx context.Context, observationID string, state model.RetryState) ([]model.RetryPlan, error) {
	items, err := listKind[model.RetryPlan](ctx, s, kindRetry)
	if err != nil {
		return nil, err
	}
	result := items[:0]
	for _, item := range items {
		if observationID != "" && item.ObservationID != observationID {
			continue
		}
		if state != "" && item.State != state {
			continue
		}
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].NextAt.Before(result[j].NextAt) })
	return result, nil
}

package store

import (
	"context"
	"fmt"
	"sort"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (s *Store) SaveAlert(ctx context.Context, alert model.Alert) error {
	if alert.ID == "" || alert.Kind == "" || alert.Message == "" {
		return fmt.Errorf("alert kind, id and message are required")
	}
	return s.Save(ctx, kindAlert, alert.ID, alert)
}

func (s *Store) GetAlert(ctx context.Context, id string) (model.Alert, error) {
	alert, err := LoadJSON[model.Alert](ctx, s, kindAlert, id)
	return alert, fmt.Errorf("load alert: %w", err)
}

func (s *Store) ListAlerts(ctx context.Context, filter model.AlertFilter) ([]model.Alert, error) {
	items, err := listKind[model.Alert](ctx, s, kindAlert)
	if err != nil {
		return nil, err
	}
	result := items[:0]
	for _, item := range items {
		if filter.ObservationID != "" && item.ObservationID != filter.ObservationID {
			continue
		}
		if filter.State != "" && item.State != filter.State {
			continue
		}
		if filter.Severity != "" && item.Severity != filter.Severity {
			continue
		}
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.After(result[j].CreatedAt) })
	if filter.Limit > 0 && len(result) > filter.Limit {
		result = result[:filter.Limit]
	}
	return result, nil
}

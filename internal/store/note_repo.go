package store

import (
	"context"
	"fmt"
	"sort"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (s *Store) SaveNote(ctx context.Context, note model.OperatorNote) error {
	if note.ID == "" || note.Operator == "" || note.Body == "" {
		return fmt.Errorf("note id, operator and body are required")
	}
	return s.Save(ctx, kindNote, note.ID, note)
}

func (s *Store) ListNotes(ctx context.Context, observationID, alertID string) ([]model.OperatorNote, error) {
	items, err := listKind[model.OperatorNote](ctx, s, kindNote)
	if err != nil {
		return nil, err
	}
	result := items[:0]
	for _, item := range items {
		if observationID != "" && item.ObservationID != observationID {
			continue
		}
		if alertID != "" && item.AlertID != alertID {
			continue
		}
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.Before(result[j].CreatedAt) })
	return result, nil
}

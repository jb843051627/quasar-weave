package store

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

const kindAudit = "audit_entry"

func (s *Store) SaveAudit(ctx context.Context, entry model.AuditEntry) error {
	if entry.ID == "" || entry.Subject == "" || entry.Action == "" || entry.OccurredAt.IsZero() {
		return fmt.Errorf("audit entry is incomplete")
	}
	return s.Save(ctx, kindAudit, entry.ID, entry)
}

func (s *Store) ListAudit(ctx context.Context, subject string, limit int) ([]model.AuditEntry, error) {
	items, err := listKind[model.AuditEntry](ctx, s, kindAudit)
	if err != nil {
		return nil, err
	}
	result := items[:0]
	for _, item := range items {
		if subject == "" || item.Subject == subject {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].OccurredAt.After(result[j].OccurredAt) })
	if limit > 0 && len(result) > limit {
		result = result[len(result)-limit:]
	}
	return result, nil
}

func (s *Store) AuditRange(ctx context.Context, start, end time.Time) ([]model.AuditEntry, error) {
	items, err := s.ListAudit(ctx, "", 0)
	if err != nil {
		return nil, err
	}
	result := items[:0]
	for _, item := range items {
		if (start.IsZero() || !item.OccurredAt.Before(start)) && (end.IsZero() || !item.OccurredAt.After(end)) {
			result = append(result, item)
		}
	}
	return result, nil
}

func (s *Store) BuildTimeline(ctx context.Context, observationID string) (model.Timeline, error) {
	events, err := s.ListAudit(ctx, observationID, 200)
	if err != nil {
		return model.Timeline{}, err
	}
	commands, err := s.ListCommands(ctx, observationID)
	if err != nil {
		return model.Timeline{}, err
	}
	return model.Timeline{ObservationID: observationID, Events: events, Commands: commands}, nil
}
